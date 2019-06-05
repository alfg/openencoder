# ENC
> FFmpeg Encoding HTTP API.

* Encoding API for submitting jobs to FFmpeg
* Redis-backed worker
* AWS S3-based storage
* Web Dashboard UI for submitting and viewing encode jobs

https://godoc.org/github.com/alfg/enc

[![Build Status](https://travis-ci.org/alfg/enc.svg?branch=master)](https://travis-ci.org/alfg/enc) 
[![GoDoc](https://godoc.org/github.com/alfg/enc?status.svg)](https://godoc.org/github.com/alfg/enc)
[![Go Report Card](https://goreportcard.com/badge/github.com/alfg/enc)](https://goreportcard.com/report/github.com/alfg/enc)

## Develop
#### Requirements
* Docker
* Go 1.11+
* FFmpeg
* Postgres
* AWS S3 Credentials & Bucket

#### Setup
* Start Redis and Postgres in Docker:
```
docker-compose up -d
```

* Start API server.
```
go build -v && enc.exe server
```

* Start worker.
```
go build -v && enc.exe worker
```

* Start Web Dashboard for development:
```
cd static && npm run serve
```

## Usage
```bash
curl -X POST \
  http://localhost:8080/api/encode \
  -H 'Content-Type: application/json' \
  -d '{
	"profile": "h264_baseline_360p_600",
	"source": "s3:///src/ToS-1080p.mp4",
	"dest": "s3:///dst/tears-of-steel/"
  }'
```

## API
TODO

## License
MIT