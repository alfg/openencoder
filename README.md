# ENC
Encoding HTTP API.

```bash
curl -X POST \
  http://localhost:8080/api/encode \
  -H 'Content-Type: application/json' \
  -d '{
	"profile": "h264_baseline_360p_600",
	"source": "s3:///src/ToS-1080p.mp4",
	"dest": "s3:///dst/tears-of-steel/",
	"delay": "2s"
  }'
```
