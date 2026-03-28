package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/example/http-stream/internal/api/httpstreamv1"
	"github.com/example/http-stream/internal/pipeline"
)

type Streamer struct {
	client   *http.Client
	registry *pipeline.Registry
	now      func() time.Time
}

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

	targetReq, err := buildHTTPRequest(ctx, req.Target, body)
	if err != nil {
		body.Close()
		return nil, fmt.Errorf("build target request: %w", err)
	}
	targetReq.ContentLength = req.Target.ContentLength

	counter := &countingReadCloser{ReadCloser: body}
	targetReq.Body = counter
	targetReq.GetBody = nil

	targetResp, err := s.client.Do(targetReq)
	if err != nil {
		counter.Close()
		return nil, fmt.Errorf("send target request: %w", err)
	}
	defer targetResp.Body.Close()
	io.Copy(io.Discard, targetResp.Body)

	return &httpstreamv1.TransferResponse{
		TransferID:       s.now().UTC().Format("20060102T150405.000000000Z07:00"),
		BytesTransferred: counter.N,
		SourceStatusCode: int32(sourceResp.StatusCode),
		TargetStatusCode: int32(targetResp.StatusCode),
	}, nil
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
	if req.Source.URL == "" || req.Target.URL == "" {
		return errors.New("source.url and target.url are required")
	}
	if req.Source.Method == "" {
		req.Source.Method = http.MethodGet
	}
	if req.Target.Method == "" {
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
