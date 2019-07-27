#### Jobs API
Jobs API resource.

| Method | Endpoint | Description |
| :----: | ---- | --------------- |
| **POST** | [/api/jobs](#create-job) | Create encode job. |
| **GET** | [/api/jobs](#list-jobs) | Get jobs list. |
| **GET** | [/api/jobs/:job_id](#get-job) | Get job details. |

---

#### Create Encode
```
POST /api/encode
```

##### Parameters
```
Content-Type: application/json
```

```json
{
    "profile": "h264_baseline_360p_600",
    "source": "s3:///src/tears-of-steel-2s.mp4",
    "dest": "s3:///dst/tears-of-steel-2s/"
}}
```

##### Response
```
Content-Type: application/json
```

```json
{
  "message": "Job created",
  "status": 200
}
```

---

#### List jobs
```
GET /api/jobs
```

##### Response
```
Content-Type: application/json
```
```json
{
  "count": 2,
  "items": [
    {
      "id": 2,
      "guid": "bkl9gbj5bidgus7kjoog",
      "profile": "h264_baseline_360p_600",
      "created_date": "2019-07-14T00:00:00Z"
    },
    {
      "id": 1,
      "guid": "bkl9gb35bidgus7kjoo0",
      "profile": "h264_baseline_360p_600",
      "created_date": "2019-07-14T00:00:00Z"
    }
  ]
}
```

---

#### Get Job
```
GET /jobs/:job_id
```

##### Response
```
Content-Type: application/json
```
```json
{
  "id": 2,
  "guid": "bkl9gbj5bidgus7kjoog",
  "profile": "h264_baseline_360p_600",
  "created_date": "2019-07-14t00:00:00z"
}
```
