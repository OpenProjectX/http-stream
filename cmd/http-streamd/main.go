package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/OpenProjectX/http-stream/internal/pipeline"
	"github.com/OpenProjectX/http-stream/internal/server"
	"github.com/OpenProjectX/http-stream/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	addr := getenv("HTTP_STREAM_LISTEN_ADDR", ":8080")

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
	streamer.SetLogger(log.Default())
	streamer.SetProgressLogInterval(getenvDuration("HTTP_STREAM_PROGRESS_LOG_INTERVAL", 2*time.Second))

	grpcServer := grpc.NewServer()
	server.Register(grpcServer, server.New(streamer))
	reflection.Register(grpcServer)

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

func getenvDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("invalid duration for %s=%q, using fallback %s", key, value, fallback)
		return fallback
	}
	return duration
}
