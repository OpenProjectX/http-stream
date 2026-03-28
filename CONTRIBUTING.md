# Contributing

This guide is written for contributors who are new to Go and want a reliable local workflow for this repository.

## Prerequisites

- Go `1.26.0`
- Git
- GoLand if you want an IDE workflow

Check your Go version:

```bash
go version
```

Expected result:

```bash
go version go1.26.0 linux/amd64
```

## Repository layout

- `cmd/http-streamd`: service entrypoint
- `internal/service`: HTTP source-to-target transfer logic
- `internal/pipeline`: streaming pipeline stages and extension points
- `internal/server`: gRPC server wiring
- `api/httpstream/v1`: proto contract

## Clone and open the project

```bash
git clone https://github.com/OpenProjectX/http-stream.git
cd http-stream
```

## Build

Build the server binary:

```bash
go build ./cmd/http-streamd
```

Build all packages:

```bash
go build ./...
```

## Run

Start the service with the default listen address:

```bash
go run ./cmd/http-streamd
```

Start the service on a custom address:

```bash
HTTP_STREAM_LISTEN_ADDR=:9090 go run ./cmd/http-streamd
```

## Test

Run the full test suite:

```bash
go test ./...
```

Run tests for one package:

```bash
go test ./internal/service
```

Run tests with verbose output:

```bash
go test -v ./...
```

## Format code

Format all Go files before committing:

```bash
gofmt -w ./cmd ./internal
```

If you changed more directories, format those files too.

## Dependency management

If you add or remove imports from external modules, update the module files:

```bash
go mod tidy
```

This updates `go.mod` and `go.sum`.

## GoLand setup

1. Open GoLand.
2. Choose `Open` and select the repository root.
3. Wait for GoLand to index the project and detect `go.mod`.
4. Make sure the project SDK points to Go `1.26.0`.
5. Open `cmd/http-streamd/main.go` if you want to run the service from the IDE.

### GoLand run configuration

To run the service from GoLand:

1. Open `Run` > `Edit Configurations`.
2. Add a new `Go Build` or `Go Application` configuration.
3. Set the package path to `github.com/OpenProjectX/http-stream/cmd/http-streamd` or select the `cmd/http-streamd` directory.
4. Optional: add `HTTP_STREAM_LISTEN_ADDR=:9090` as an environment variable.
5. Run the configuration.

### GoLand test workflow

- Open a `_test.go` file and click the gutter run icon.
- Or right-click a package directory such as `internal/service` and choose `Run 'Go Tests in ...'`.

## How to add a new pipeline stage

1. Create a new type in `internal/pipeline`.
2. Make it satisfy the `pipeline.Stage` interface.
3. Register it in `cmd/http-streamd/main.go`.
4. Add unit tests for the stage.
5. Update `README.md` if the new stage is user-facing.

Example extension points:

- compression
- encryption
- checksumming
- throttling

## Suggested contribution workflow

1. Create a branch.
2. Make a focused change.
3. Run `gofmt -w ./cmd ./internal`.
4. Run `go test ./...`.
5. Commit with a clear message.
6. Open a pull request.

## Common beginner notes

- `internal/...` packages are intentionally private to this module.
- Streaming code usually works with `io.Reader`, `io.Writer`, and `io.ReadCloser`.
- Avoid loading the entire payload into memory unless there is a strong reason.
- Prefer small, testable packages over large files with mixed responsibilities.

## Before opening a pull request

Please make sure:

- the project builds
- tests pass
- code is formatted
- new behavior has tests where practical
- docs are updated when the public behavior changes
