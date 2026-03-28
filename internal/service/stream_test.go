package service

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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
}
