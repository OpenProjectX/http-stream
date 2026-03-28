package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/OpenProjectX/http-stream/internal/api/httpstreamv1"
	"github.com/OpenProjectX/http-stream/internal/pipeline"
)

type Streamer struct {
	client   *http.Client
	registry *pipeline.Registry
	now      func() time.Time
}

type ProgressObserver func(*httpstreamv1.TransferProgress) error

func New(client *http.Client, registry *pipeline.Registry) *Streamer {
	if client == nil {
		client = &http.Client{Timeout: 0}
	}
	if registry == nil {
		registry = pipeline.NewRegistry()
	}
	return &Streamer{
		client:   client,
		registry: registry,
		now:      time.Now,
	}
}

func (s *Streamer) Transfer(ctx context.Context, req *httpstreamv1.TransferRequest) (*httpstreamv1.TransferResponse, error) {
	resp, err := s.transfer(ctx, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Streamer) TransferStream(ctx context.Context, req *httpstreamv1.TransferRequest, observer ProgressObserver) error {
	_, err := s.transfer(ctx, req, observer)
	return err
}

func (s *Streamer) transfer(ctx context.Context, req *httpstreamv1.TransferRequest, observer ProgressObserver) (*httpstreamv1.TransferResponse, error) {
	startedAt := s.now()
	transferID := startedAt.UTC().Format("20060102T150405.000000000Z07:00")

	if err := validateTransferRequest(req); err != nil {
		return nil, err
	}

	sourceReq, err := buildHTTPRequest(ctx, req.Source, nil)
	if err != nil {
		return nil, fmt.Errorf("build source request: %w", err)
	}

	sourceResp, err := s.client.Do(sourceReq)
	if err != nil {
		return nil, fmt.Errorf("send source request: %w", err)
	}

	if sourceResp.StatusCode < http.StatusOK || sourceResp.StatusCode >= http.StatusMultipleChoices {
		defer sourceResp.Body.Close()
		return nil, fmt.Errorf("source request failed with status %d", sourceResp.StatusCode)
	}

	if err := emitProgress(observer, buildTransferProgress(transferID, s.now, startedAt, 0, sourceResp.ContentLength, int32(sourceResp.StatusCode), 0, false)); err != nil {
		sourceResp.Body.Close()
		return nil, err
	}

	stageSpecs := make([]pipeline.StageSpec, 0, len(req.Pipeline))
	for _, stage := range req.Pipeline {
		stageSpecs = append(stageSpecs, pipeline.StageSpec{
			Name:   stage.Name,
			Config: pipeline.StageConfig(stage.Config),
		})
	}

	body, err := s.registry.Build(ctx, sourceResp.Body, stageSpecs)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	progressReader := &progressReadCloser{
		ReadCloser:          body,
		transferID:          transferID,
		now:                 s.now,
		startedAt:           startedAt,
		sourceContentLength: sourceResp.ContentLength,
		sourceStatusCode:    int32(sourceResp.StatusCode),
		observer:            observer,
	}

	if req.Target.LocalPath != "" {
		return s.transferToLocalFile(progressReader, sourceResp.StatusCode, sourceResp.ContentLength, req.Target.LocalPath, transferID, startedAt)
	}

	targetReq, err := buildHTTPRequest(ctx, req.Target, progressReader)
	if err != nil {
		return nil, fmt.Errorf("build target request: %w", err)
	}
	targetReq.ContentLength = req.Target.ContentLength

	targetReq.Body = progressReader
	targetReq.GetBody = nil

	targetResp, err := s.client.Do(targetReq)
	if err != nil {
		return nil, fmt.Errorf("send target request: %w", err)
	}
	defer targetResp.Body.Close()
	io.Copy(io.Discard, targetResp.Body)
	if targetResp.StatusCode < http.StatusOK || targetResp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("target request failed with status %d", targetResp.StatusCode)
	}

	resp := buildTransferResponse(transferID, s.now, startedAt, progressReader.N, sourceResp.ContentLength, int32(sourceResp.StatusCode), int32(targetResp.StatusCode))
	if err := emitProgress(observer, progressFromResponse(resp, true)); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Streamer) transferToLocalFile(body io.Reader, sourceStatusCode int, sourceContentLength int64, localPath, transferID string, startedAt time.Time) (*httpstreamv1.TransferResponse, error) {
	if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
		return nil, fmt.Errorf("create parent directories for %q: %w", localPath, err)
	}

	file, err := os.Create(localPath)
	if err != nil {
		return nil, fmt.Errorf("create target file %q: %w", localPath, err)
	}
	defer file.Close()

	written, err := io.Copy(file, body)
	if err != nil {
		return nil, fmt.Errorf("write target file %q: %w", localPath, err)
	}

	resp := buildTransferResponse(transferID, s.now, startedAt, written, sourceContentLength, int32(sourceStatusCode), 0)
	if progressReader, ok := body.(*progressReadCloser); ok {
		progressReader.N = written
		if err := emitProgress(progressReader.observer, progressFromResponse(resp, true)); err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func buildTransferResponse(transferID string, now func() time.Time, startedAt time.Time, bytesTransferred, sourceContentLength int64, sourceStatusCode, targetStatusCode int32) *httpstreamv1.TransferResponse {
	finishedAt := now()
	duration := finishedAt.Sub(startedAt)
	durationMillis := duration.Milliseconds()
	if durationMillis < 0 {
		durationMillis = 0
	}

	var averageBytesPerSecond float64
	if duration > 0 {
		averageBytesPerSecond = float64(bytesTransferred) / duration.Seconds()
	}

	progressPercent := 100.0
	if sourceContentLength > 0 {
		progressPercent = (float64(bytesTransferred) / float64(sourceContentLength)) * 100
		if progressPercent > 100 {
			progressPercent = 100
		}
	}

	return &httpstreamv1.TransferResponse{
		TransferID:            transferID,
		BytesTransferred:      bytesTransferred,
		SourceStatusCode:      sourceStatusCode,
		TargetStatusCode:      targetStatusCode,
		SourceContentLength:   sourceContentLength,
		DurationMillis:        durationMillis,
		AverageBytesPerSecond: averageBytesPerSecond,
		ProgressPercent:       progressPercent,
	}
}

func buildTransferProgress(transferID string, now func() time.Time, startedAt time.Time, bytesTransferred, sourceContentLength int64, sourceStatusCode, targetStatusCode int32, done bool) *httpstreamv1.TransferProgress {
	finishedAt := now()
	duration := finishedAt.Sub(startedAt)
	durationMillis := duration.Milliseconds()
	if durationMillis < 0 {
		durationMillis = 0
	}

	var averageBytesPerSecond float64
	if duration > 0 {
		averageBytesPerSecond = float64(bytesTransferred) / duration.Seconds()
	}

	progressPercent := 0.0
	if bytesTransferred > 0 {
		progressPercent = 100.0
	}
	if sourceContentLength > 0 {
		progressPercent = (float64(bytesTransferred) / float64(sourceContentLength)) * 100
		if progressPercent > 100 {
			progressPercent = 100
		}
	}
	if done && progressPercent == 0 && sourceContentLength <= 0 {
		progressPercent = 100
	}

	return &httpstreamv1.TransferProgress{
		TransferID:            transferID,
		BytesTransferred:      bytesTransferred,
		SourceStatusCode:      sourceStatusCode,
		TargetStatusCode:      targetStatusCode,
		SourceContentLength:   sourceContentLength,
		DurationMillis:        durationMillis,
		AverageBytesPerSecond: averageBytesPerSecond,
		ProgressPercent:       progressPercent,
		Done:                  done,
	}
}

func progressFromResponse(resp *httpstreamv1.TransferResponse, done bool) *httpstreamv1.TransferProgress {
	return &httpstreamv1.TransferProgress{
		TransferID:            resp.TransferID,
		BytesTransferred:      resp.BytesTransferred,
		SourceStatusCode:      resp.SourceStatusCode,
		TargetStatusCode:      resp.TargetStatusCode,
		SourceContentLength:   resp.SourceContentLength,
		DurationMillis:        resp.DurationMillis,
		AverageBytesPerSecond: resp.AverageBytesPerSecond,
		ProgressPercent:       resp.ProgressPercent,
		Done:                  done,
	}
}

func emitProgress(observer ProgressObserver, progress *httpstreamv1.TransferProgress) error {
	if observer == nil {
		return nil
	}
	return observer(progress)
}

func validateTransferRequest(req *httpstreamv1.TransferRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	if req.Source == nil {
		return errors.New("source is required")
	}
	if req.Target == nil {
		return errors.New("target is required")
	}
	if req.Source.URL == "" {
		return errors.New("source.url is required")
	}
	if req.Target.URL == "" && req.Target.LocalPath == "" {
		return errors.New("target.url or target.local_path is required")
	}
	if req.Target.URL != "" && req.Target.LocalPath != "" {
		return errors.New("target.url and target.local_path are mutually exclusive")
	}
	if req.Source.Method == "" {
		req.Source.Method = http.MethodGet
	}
	if req.Target.LocalPath == "" && req.Target.Method == "" {
		req.Target.Method = http.MethodPut
	}
	return nil
}

func buildHTTPRequest(ctx context.Context, spec *httpstreamv1.HTTPRequest, body io.ReadCloser) (*http.Request, error) {
	var reader io.Reader
	if body != nil {
		reader = body
	}

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(spec.Method), spec.URL, reader)
	if err != nil {
		if body != nil {
			body.Close()
		}
		return nil, err
	}

	for k, v := range spec.Headers {
		req.Header.Set(k, v)
	}

	if body != nil {
		req.Body = body
	}

	return req, nil
}

type countingReadCloser struct {
	io.ReadCloser
	N int64
}

func (c *countingReadCloser) Read(p []byte) (int, error) {
	n, err := c.ReadCloser.Read(p)
	c.N += int64(n)
	return n, err
}

type progressReadCloser struct {
	io.ReadCloser
	N                   int64
	transferID          string
	now                 func() time.Time
	startedAt           time.Time
	sourceContentLength int64
	sourceStatusCode    int32
	observer            ProgressObserver
	pendingErr          error
}

func (p *progressReadCloser) Read(buf []byte) (int, error) {
	if p.pendingErr != nil {
		return 0, p.pendingErr
	}

	n, err := p.ReadCloser.Read(buf)
	p.N += int64(n)
	if n > 0 {
		progress := buildTransferProgress(p.transferID, p.now, p.startedAt, p.N, p.sourceContentLength, p.sourceStatusCode, 0, false)
		if notifyErr := emitProgress(p.observer, progress); notifyErr != nil {
			p.pendingErr = notifyErr
			if err == nil {
				err = notifyErr
			}
		}
	}
	return n, err
}
