package service

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/OpenProjectX/http-stream/internal/api/httpstreamv1"
	"github.com/OpenProjectX/http-stream/internal/pipeline"
)

func TestTransfer(t *testing.T) {
	payload := "stream me"
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("source method = %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, payload)
	}))
	defer source.Close()

	var received string
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		received = string(body)
		w.WriteHeader(http.StatusCreated)
	}))
	defer target.Close()

	svc := New(source.Client(), pipeline.NewRegistry())
	svc.now = func() time.Time {
		return time.Date(2026, 3, 28, 13, 0, 0, 0, time.UTC)
	}
	resp, err := svc.Transfer(context.Background(), &httpstreamv1.TransferRequest{
		Source: &httpstreamv1.HTTPRequest{
			Method: http.MethodGet,
			URL:    source.URL,
		},
		Target: &httpstreamv1.HTTPRequest{
			Method: http.MethodPut,
			URL:    target.URL,
		},
	})
	if err != nil {
		t.Fatalf("Transfer() error = %v", err)
	}

	if received != payload {
		t.Fatalf("received = %q want %q", received, payload)
	}
	if resp.BytesTransferred != int64(len(payload)) {
		t.Fatalf("bytes = %d want %d", resp.BytesTransferred, len(payload))
	}
	if resp.SourceStatusCode != http.StatusOK {
		t.Fatalf("source status = %d", resp.SourceStatusCode)
	}
	if resp.TargetStatusCode != http.StatusCreated {
		t.Fatalf("target status = %d", resp.TargetStatusCode)
	}
	if resp.SourceContentLength != int64(len(payload)) {
		t.Fatalf("source content length = %d want %d", resp.SourceContentLength, len(payload))
	}
	if resp.ProgressPercent != 100 {
		t.Fatalf("progress = %f want 100", resp.ProgressPercent)
	}
}

func TestTransferToLocalFile(t *testing.T) {
	payload := "disk target"
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, payload)
	}))
	defer source.Close()

	targetPath := filepath.Join(t.TempDir(), "nested", "payload.bin")

	svc := New(source.Client(), pipeline.NewRegistry())
	svc.now = func() time.Time {
		return time.Date(2026, 3, 28, 13, 0, 0, 0, time.UTC)
	}
	resp, err := svc.Transfer(context.Background(), &httpstreamv1.TransferRequest{
		Source: &httpstreamv1.HTTPRequest{
			Method: http.MethodGet,
			URL:    source.URL,
		},
		Target: &httpstreamv1.HTTPRequest{
			LocalPath: targetPath,
		},
	})
	if err != nil {
		t.Fatalf("Transfer() error = %v", err)
	}

	got, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(got) != payload {
		t.Fatalf("file contents = %q want %q", got, payload)
	}
	if resp.BytesTransferred != int64(len(payload)) {
		t.Fatalf("bytes = %d want %d", resp.BytesTransferred, len(payload))
	}
	if resp.TargetStatusCode != 0 {
		t.Fatalf("target status = %d want 0 for local file target", resp.TargetStatusCode)
	}
	if resp.SourceContentLength != int64(len(payload)) {
		t.Fatalf("source content length = %d want %d", resp.SourceContentLength, len(payload))
	}
	if resp.ProgressPercent != 100 {
		t.Fatalf("progress = %f want 100", resp.ProgressPercent)
	}
}

func TestTransferStreamProgress(t *testing.T) {
	payload := "progress payload"
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, payload)
	}))
	defer source.Close()

	targetPath := filepath.Join(t.TempDir(), "progress", "payload.bin")

	svc := New(source.Client(), pipeline.NewRegistry())
	base := time.Date(2026, 3, 28, 13, 0, 0, 0, time.UTC)
	calls := 0
	svc.now = func() time.Time {
		defer func() { calls++ }()
		return base.Add(time.Duration(calls) * 100 * time.Millisecond)
	}

	var updates []*httpstreamv1.TransferProgress
	err := svc.TransferStream(context.Background(), &httpstreamv1.TransferRequest{
		Source: &httpstreamv1.HTTPRequest{
			Method: http.MethodGet,
			URL:    source.URL,
		},
		Target: &httpstreamv1.HTTPRequest{
			LocalPath: targetPath,
		},
	}, func(progress *httpstreamv1.TransferProgress) error {
		copy := *progress
		updates = append(updates, &copy)
		return nil
	})
	if err != nil {
		t.Fatalf("TransferStream() error = %v", err)
	}

	if len(updates) < 2 {
		t.Fatalf("updates len = %d want at least 2", len(updates))
	}
	if updates[0].BytesTransferred != 0 || updates[0].Done {
		t.Fatalf("first update = %+v want initial progress event", updates[0])
	}
	last := updates[len(updates)-1]
	if !last.Done {
		t.Fatalf("last update done = %v want true", last.Done)
	}
	if last.BytesTransferred != int64(len(payload)) {
		t.Fatalf("last bytes = %d want %d", last.BytesTransferred, len(payload))
	}
	if last.ProgressPercent != 100 {
		t.Fatalf("last progress = %f want 100", last.ProgressPercent)
	}
}

func TestProgressReadCloserEmitsOnSizeWindow(t *testing.T) {
	base := time.Date(2026, 3, 28, 13, 0, 0, 0, time.UTC)
	times := []time.Time{
		base.Add(100 * time.Millisecond),
		base.Add(200 * time.Millisecond),
		base.Add(200 * time.Millisecond),
		base.Add(300 * time.Millisecond),
	}
	calls := 0

	reader := &progressReadCloser{
		ReadCloser:          io.NopCloser(&fixedChunkReader{chunks: [][]byte{[]byte("ab"), []byte("cd"), []byte("ef")}}),
		transferID:          "size-window",
		now:                 func() time.Time { now := times[calls]; calls++; return now },
		startedAt:           base,
		sourceContentLength: 6,
		sourceStatusCode:    http.StatusOK,
		logger:              log.New(io.Discard, "", 0),
		logInterval:         time.Second,
		emitInterval:        time.Second,
		emitBytes:           4,
		lastEmitAt:          base,
	}

	var updates []*httpstreamv1.TransferProgress
	reader.observer = func(progress *httpstreamv1.TransferProgress) error {
		copy := *progress
		updates = append(updates, &copy)
		return nil
	}

	buf := make([]byte, 2)
	for {
		_, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
	}

	if len(updates) != 1 {
		t.Fatalf("updates len = %d want 1", len(updates))
	}
	if updates[0].BytesTransferred != 4 {
		t.Fatalf("bytes transferred = %d want 4", updates[0].BytesTransferred)
	}
}

func TestProgressReadCloserEmitsOnTimeWindow(t *testing.T) {
	base := time.Date(2026, 3, 28, 13, 0, 0, 0, time.UTC)
	times := []time.Time{
		base.Add(400 * time.Millisecond),
		base.Add(1100 * time.Millisecond),
		base.Add(1100 * time.Millisecond),
		base.Add(1500 * time.Millisecond),
	}
	calls := 0

	reader := &progressReadCloser{
		ReadCloser:          io.NopCloser(&fixedChunkReader{chunks: [][]byte{[]byte("ab"), []byte("cd"), []byte("ef")}}),
		transferID:          "time-window",
		now:                 func() time.Time { now := times[calls]; calls++; return now },
		startedAt:           base,
		sourceContentLength: 6,
		sourceStatusCode:    http.StatusOK,
		logger:              log.New(io.Discard, "", 0),
		logInterval:         time.Second,
		emitInterval:        time.Second,
		emitBytes:           10,
		lastEmitAt:          base,
	}

	var updates []*httpstreamv1.TransferProgress
	reader.observer = func(progress *httpstreamv1.TransferProgress) error {
		copy := *progress
		updates = append(updates, &copy)
		return nil
	}

	buf := make([]byte, 2)
	for {
		_, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
	}

	if len(updates) != 1 {
		t.Fatalf("updates len = %d want 1", len(updates))
	}
	if updates[0].BytesTransferred != 4 {
		t.Fatalf("bytes transferred = %d want 4", updates[0].BytesTransferred)
	}
}

type fixedChunkReader struct {
	chunks [][]byte
	index  int
}

func (r *fixedChunkReader) Read(p []byte) (int, error) {
	if r.index >= len(r.chunks) {
		return 0, io.EOF
	}
	n := copy(p, r.chunks[r.index])
	r.index++
	return n, nil
}
