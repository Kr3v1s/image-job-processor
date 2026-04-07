# Image Job Processor

Concurrent backend service built in Go that demonstrates worker pools, async processing, and fault tolerance.

## Overview

This project exposes a small HTTP API that accepts image processing jobs, downloads the source image, generates a thumbnail, and lets clients poll for job status.

It is intentionally simple:

- in-memory job storage
- background processing with workers
- thumbnail generation on disk
- no database, auth, or framework

The goal is to demonstrate backend fundamentals that matter in production-oriented Go services:

- concurrency with goroutines and channels
- worker pool design
- thread-safe shared state with `sync.RWMutex`
- resilient HTTP download flow with timeout and retry
- clear async API behavior with job polling

## Architecture

```text
Client -> HTTP API -> In-Memory Job Store -> Queue -> Workers -> Download + Thumbnail
```

## Processing Flow

1. Client submits a job with `POST /jobs`.
2. The API creates a job in `pending` state and stores it in memory.
3. The job ID is pushed onto a buffered queue.
4. A worker picks up the job and marks it as `running`.
5. The worker downloads the image with timeout and retry logic.
6. The original image is saved to `output/`.
7. A thumbnail is generated and saved to `output/`.
8. The job is updated to `done` or `failed`.
9. Client checks progress with `GET /jobs?id=...`.

## API

### `POST /jobs`

Creates a new image processing job.

Request:

```http
POST /jobs
Content-Type: application/json
```

```json
{
  "image_url": "https://example.com/image.jpg"
}
```

Example response:

```json
{
  "id": "1744100000000000000",
  "status": "pending",
  "image_url": "https://example.com/image.jpg",
  "result": "",
  "thumbnail": "",
  "error": ""
}
```

### `GET /jobs?id=<job_id>`

Returns the current status of a job.

Completed job example:

```json
{
  "id": "1744100000000000000",
  "status": "done",
  "image_url": "https://example.com/image.jpg",
  "result": "output/1744100000000000000.jpg",
  "thumbnail": "output/thumb_1744100000000000000.jpg",
  "error": ""
}
```

Failed job example:

```json
{
  "id": "1744100000000000000",
  "status": "failed",
  "image_url": "https://example.com/image.jpg",
  "result": "",
  "thumbnail": "",
  "error": "failed after retries: status 404 Not Found"
}
```

## Project Structure

```text
.
|-- main.go
|-- internal/
|   |-- http/
|   |   `-- handler.go
|   `-- jobs/
|       |-- image.go
|       |-- job.go
|       `-- worker.go
|-- output/
|-- go.mod
`-- README.md
```

## Key Implementation Details

### Worker Pool

Multiple workers consume job IDs from a buffered channel, which keeps the HTTP layer fast and the processing flow asynchronous.

### Thread-Safe Job Store

Jobs are stored in memory behind a store protected by `sync.RWMutex`, avoiding race conditions between HTTP handlers and background workers.

### Resilient Downloads

Workers use an HTTP client with timeout and retry logic so temporary network failures do not immediately fail the job.

### File Outputs

Each successful job produces:

- the original downloaded image in `output/`
- a generated thumbnail in `output/`

## Run Locally

Requirements:

- Go 1.22+ installed

Run:

```bash
go mod tidy
go run main.go
```

The server starts on `http://localhost:8080`.

Note: the project writes processed files into the `output/` directory.

## Quick Test

Create a job:

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d "{\"image_url\":\"https://via.placeholder.com/600\"}"
```

Check job status:

```bash
curl "http://localhost:8080/jobs?id=<job_id>"
```

Health check:

```bash
curl http://localhost:8080/health
```

## Current Tradeoffs

This project is intentionally scoped for speed and clarity:

- jobs are stored in memory only
- processed files are written to local disk
- there is no authentication or rate limiting
- job status is retrieved by polling

Those tradeoffs keep the code easy to understand while still showing solid backend fundamentals.

## Next Improvements

- add structured logging
- add request logging middleware
- add tests for handlers and worker behavior
- add Docker support
- add persistent storage for jobs
- expose basic metrics

## Why This Project

This project was built as a portfolio-ready backend exercise focused on practical Go engineering:

- asynchronous processing
- concurrency control
- error handling
- API design
- operationally aware code without unnecessary architecture

## Author

Kevin Bustamante  
Backend engineer focused on building reliable systems with Go.
