package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/example/http-stream/internal/pipeline"
	"github.com/example/http-stream/internal/server"
	"github.com/example/http-stream/internal/service"
	"github.com/example/http-stream/internal/transport/grpcjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

func main() {
	addr := getenv("HTTP_STREAM_LISTEN_ADDR", ":8080")

	encoding.RegisterCodec(grpcjson.Codec{})

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}

	httpClient := &http.Client{
		Timeout: 0,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	registry := pipeline.NewRegistry(pipeline.AESCTRStage{})
	streamer := service.New(httpClient, registry)

	grpcServer := grpc.NewServer()
	server.Register(grpcServer, server.New(streamer))

	log.Printf("http-streamd listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve grpc: %v", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
