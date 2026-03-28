# http-stream

`http-stream` is a Go 1.26 service that copies an HTTP source response body directly into either an HTTP target request body or a local file without buffering the full payload to disk. The service is controlled through a gRPC API and supports pluggable streaming pipeline stages, so transforms such as encryption can be applied in flight.

## What is implemented

- gRPC service contract in [`api/httpstream/v1/httpstream.proto`](api/httpstream/v1/httpstream.proto)
- gRPC server in [`cmd/http-streamd`](cmd/http-streamd)
- HTTP source-to-target streaming with zero-disk buffering
- HTTP source-to-local-disk streaming
- Extensible pipeline registry over `io.Reader` / `io.ReadCloser`
- Built-in `encrypt.aes_ctr` stage for streaming encryption before upload
- Unit tests for the pipeline stage and the HTTP transfer flow

## API shape

The service exposes one RPC:

```proto
rpc Transfer(TransferRequest) returns (TransferResponse);
rpc TransferStream(TransferRequest) returns (stream TransferProgress);
```

`TransferRequest` contains:

- `source`: upstream HTTP request definition, typically `GET`
- `target`: downstream HTTP request definition for remote uploads, or a local disk path for on-machine writes
- `pipeline`: ordered stages to wrap the source stream before it is sent to the target

`TransferResponse` now includes transfer observability fields:

- `bytesTransferred`: final streamed byte count
- `sourceContentLength`: source response content length when known
- `durationMillis`: end-to-end transfer duration inside the service
- `averageBytesPerSecond`: average throughput over the completed transfer
- `progressPercent`: final completion percentage, typically `100` for successful transfers

`TransferStream` is the UI-oriented variant. It emits:

- an initial event with `bytesTransferred = 0`
- intermediate progress events while bytes are flowing
- a final event with `done = true`

## Example request

The service now speaks standard protobuf gRPC, so IDE clients such as GoLand or IntelliJ gRPC requests can use the proto contract directly.

HTTP target example:

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

Local disk target example:

```json
{
  "source": {
    "method": "GET",
    "url": "https://source.example/object"
  },
  "target": {
    "localPath": "/tmp/http-stream/object.bin"
  }
}
```

## Run

```bash
go run ./cmd/http-streamd
```

Set `HTTP_STREAM_LISTEN_ADDR` to override the default listen address `:8080`.
Set `HTTP_STREAM_PROGRESS_LOG_INTERVAL` to control periodic progress logs. Example: `500ms`, `2s`, `5s`.

## Test

```bash
go test ./...
```

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for setup instructions, GoLand guidance, and the usual build, run, format, and test workflow.

## Extension model

To add a new pipeline transform:

1. Implement `pipeline.Stage`.
2. Register it in the registry used by the server.
3. Pass the stage name and config in the `pipeline` field of the request.

This keeps transport concerns, HTTP streaming, and transform logic separated so the project can grow into retries, observability, auth plugins, or richer transfer policies without reworking the core pipeline.

For local disk targets, the service creates parent directories automatically and returns `target_status_code = 0` because no downstream HTTP response exists.

For observability, the gRPC response now carries final transfer metrics so clients can record throughput and completion without scraping logs.

For continuous observability, use `TransferStream` from the UI and bind `progressPercent` or `bytesTransferred / sourceContentLength` to the progress bar.

## Troubleshooting Logs

The service now logs detailed transfer lifecycle events for debugging:

- transfer start with source, target, and pipeline size
- source response status and content length
- target selection and request failures
- periodic progress snapshots with rate and percent complete
- final completion metrics

Progress logging is rate-limited by `HTTP_STREAM_PROGRESS_LOG_INTERVAL` and defaults to `2s`.
