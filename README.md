# http-stream

`http-stream` is a Go 1.26 service that copies an HTTP source response body directly into an HTTP target request body without buffering the full payload to disk. The service is controlled through a gRPC API and supports pluggable streaming pipeline stages, so transforms such as encryption can be applied in flight.

## What is implemented

- gRPC service contract in [`api/httpstream/v1/httpstream.proto`](api/httpstream/v1/httpstream.proto)
- gRPC server in [`cmd/http-streamd`](cmd/http-streamd)
- HTTP source-to-target streaming with zero-disk buffering
- Extensible pipeline registry over `io.Reader` / `io.ReadCloser`
- Built-in `encrypt.aes_ctr` stage for streaming encryption before upload
- Unit tests for the pipeline stage and the HTTP transfer flow

## API shape

The service exposes one RPC:

```proto
rpc Transfer(TransferRequest) returns (TransferResponse);
```

`TransferRequest` contains:

- `source`: upstream HTTP request definition, typically `GET`
- `target`: downstream HTTP request definition, typically `PUT` or `POST`
- `pipeline`: ordered stages to wrap the source stream before it is sent to the target

## Example request

The runtime uses a JSON gRPC codec for now, which keeps the service buildable without requiring `protoc` in the environment. The proto file remains the contract to generate standard stubs later.

```json
{
  "source": {
    "method": "GET",
    "url": "https://source.example/object",
    "headers": {
      "Authorization": "Bearer source-token"
    }
  },
  "target": {
    "method": "PUT",
    "url": "https://target.example/object",
    "headers": {
      "Authorization": "Bearer target-token",
      "Content-Type": "application/octet-stream"
    }
  },
  "pipeline": [
    {
      "name": "encrypt.aes_ctr",
      "config": {
        "key_b64": "<base64-encoded-16-24-or-32-byte-key>",
        "iv_b64": "<base64-encoded-16-byte-iv>"
      }
    }
  ]
}
```

## Run

```bash
go run ./cmd/http-streamd
```

Set `HTTP_STREAM_LISTEN_ADDR` to override the default listen address `:8080`.

## Test

```bash
go test ./...
```

## Extension model

To add a new pipeline transform:

1. Implement `pipeline.Stage`.
2. Register it in the registry used by the server.
3. Pass the stage name and config in the `pipeline` field of the request.

This keeps transport concerns, HTTP streaming, and transform logic separated so the project can grow into retries, observability, auth plugins, or richer transfer policies without reworking the core pipeline.
