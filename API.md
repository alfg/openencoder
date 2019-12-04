## Open Encoder API v1

### Overview
Early draft of the Open Encoder API. This is in early development and subject to change without notice.

### Resources

* [Authentication](#authentication)
* [Jobs](#jobs)
* [Machines](#machines)
* [Presets](#presets)

---

#### Authentication
Authentication API resource.

| Method | Endpoint | Description |
| :----: | ---- | --------------- |
| **POST** | [/api/register](#register) | Register a user. |
| **POST** | [/api/login](#login) | Login user. |

---

#### Register
```
POST /api/register
```

##### Parameters
```
Content-Type: application/json
```

```json
{
    "username": "foo@bar.com",
    "password": "baz"
}
```

##### Response
```
Content-Type: application/json
```

```json
{
  "user": "foo@bar.com",
  "message": "User registered"
}
```

#### Login
```
POST /api/login
```

##### Parameters
```
Content-Type: application/json
```

```json
{
  "username": "foo@bar.com",
  "password": "baz"
}
```

##### Response
```
Content-Type: application/json
```

```json
{
  "token": "<jwt token>"
}
```

#### Jobs
Jobs API resource.

| Method | Endpoint | Description |
| :----: | ---- | --------------- |
| **POST** | [/api/jobs](#create-job) | Create encode job. |
| **GET** | [/api/jobs](#list-jobs) | Get jobs list. |
| **GET** | [/api/jobs/:job_id](#get-job) | Get job details. |
| **GET** | [/api/jobs/:job_id/status](#get-job-status) | Get job status. |
| **POST** | [/api/jobs/:job_id/cancel](#cancel-job) | Cancel job. |
| **POST** | [/api/jobs/:job_id/restart](#restart-job) | Restart job. |


---

#### Create Job
```
POST /api/job
```

##### Parameters
```
Content-Type: application/json
```

```json
{
    "preset": "h264_baseline_360p_600",
    "source": "s3:///src/tears-of-steel-2s.mp4",
    "dest": "s3:///dst/tears-of-steel-2s/"
}
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
      "preset": "h264_baseline_360p_600",
      "created_date": "2019-07-14T00:00:00Z"
    },
    {
      "id": 1,
      "guid": "bkl9gb35bidgus7kjoo0",
      "preset": "h264_baseline_360p_600",
      "created_date": "2019-07-14T00:00:00Z"
    }
  ]
}
```

---

#### Get Job
```
GET /api/jobs/:job_id
```

##### Response
```
Content-Type: application/json
```
```json
{
  "id": 2,
  "guid": "bkl9gbj5bidgus7kjoog",
  "preset": "h264_baseline_360p_600",
  "created_date": "2019-07-14t00:00:00z"
}
```

---

#### Get Job Status
```
GET /api/jobs/:job_id/status
```

##### Response
```
Content-Type: application/json
```
```json
{
  "status": 200,
  "job_status": "encoding"
}
```

---

#### Cancel Job
```
POST /api/jobs/:job_id/cancel
```

##### Response
```
Content-Type: application/json
```
```json
{
  "status": 200,
}
```

---

#### Restart Job
```
POST /api/jobs/:job_id/restart
```

##### Response
```
Content-Type: application/json
```
```json
{
  "status": 200,
}
```

#### Machines
Machines API resource.

| Method | Endpoint | Description |
| :----: | ---- | --------------- |
| **POST** | [/api/machines](#create-machine) | Create a machine. |
| **GET** | [/api/machines](#list-machines) | Get machines list. |
| **DELETE** | [/api/machines/:id](#delete-machine) | Delete a machine by ID. |
| **DELETE** | [/api/machines](#delete-all-machines) | Delete all machines by tag. |
| **GET** | [/api/machines/regions](#list-machine-regions) | Get machine regions list. |
| **GET** | [/api/machines/sizes](#list-machine-sizes) | Get machine sizes list. |

---

#### Create Machine
```
POST /api/machines
```

##### Parameters
```
Content-Type: application/json
```

```json
{
	"provider": "digitalocean",
	"region": "sfo1",
	"size": "s-1vcpu-1gb",
	"count": 1
}
```

##### Response
```
Content-Type: application/json
```

```json
{
  "machine": [
    {
      "id": 154372950,
      "provider": "digitalocean"
    }
  ]
}
```

---

#### List Machines
```
GET /api/machines
```

##### Response
```
Content-Type: application/json
```
```json
{
  "machines": [
    {
      "id": 154485626,
      "name": "openencoder-worker",
      "status": "new",
      "size_slug": "s-1vcpu-1gb",
      "created_at": "2019-08-11T05:50:26Z",
      "region": "San Francisco 1",
      "tags": [
        "openencoder"
      ],
      "provider": "digitalocean"
    }
  ]
}
```

---

#### Delete Machine
```
DELETE /api/machines/:id
```

##### Parameters
```
Content-Type: application/json
```

##### Response
```
Content-Type: application/json
```

```json
{
  "machine": {
    "id": 154472730,
    "provider": "digitalocean"
  }
}
```

---

#### Delete All Machines
This will delete all machines tagged with `openencoder`.

```
DELETE /api/machine
```

##### Parameters
```
Content-Type: application/json
```

##### Response
```
Content-Type: application/json
```

```json
{
  "deleted": true
}
```

---

#### List Machine Regions
```
GET /api/machines/regions
```

##### Response
```
Content-Type: application/json
```
```json
{
  "regions": [
    {
      "name": "New York 1",
      "slug": "nyc1",
      "sizes": [
        "32gb",
        "16gb"
      ],
      "available": true
    }
  ]
}
```

---

#### List Machine Sizes
```
GET /api/machines/sizes
```

##### Response
```
Content-Type: application/json
```
```json
{
  "sizes": [
    {
      "slug": "512mb",
      "available": true,
      "price_monthly": 5,
      "price_hourly": 0.007439999841153622
    },
    {
      "slug": "s-1vcpu-1gb",
      "available": true,
      "price_monthly": 5,
      "price_hourly": 0.007439999841153622
    }
  ]
}
```

---

#### Presets
Presets API resource.

| Method | Endpoint | Description |
| :----: | ---- | --------------- |
| **POST** | [/api/presets](#create-preset) | Create preset. |
| **GET** | [/api/presets](#list-presets) | Get presets list. |
| **GET** | [/api/presets/:preset_id](#get-preset) | Get preset details. |
| **PUT** | [/api/presets/:preset_id](#update-preset) | Update preset. |

---

#### Create Preset
```
POST /api/presets
```

##### Parameters
```
Content-Type: application/json
```

```json
{
    "name": "preset name",
    "description": "preset description",
    "data": "<json string of ffmpeg preset data>",
    "active": true
}
```

##### Response
```
Content-Type: application/json
```

```json
{
  "id": 1,
  "name": "preset name",
  "description": "preset description",
  "data": "<json string of ffmpeg preset data>",
  "active": true
}
```

---

#### List Presets
```
GET /api/presets
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
      "id": 4,
      "name": "h264_main_1080p_6000",
      "description": "h264_main_1080p_6000",
      "data": "<json data>",
      "active": true
    },
    {
      "id": 3,
      "name": "h264_main_720p_3000",
      "description": "h264_main_720p_3000",
      "data": "<json data>",
      "active": true
    },
  ]
}
```

---

#### Get Preset
```
GET /api/presets/:preset_id
```

##### Response
```
Content-Type: application/json
```
```json
{
  "preset": {
    "id": 1,
    "name": "h264_baseline_360p_600",
    "description": "h264_baseline_360p_600",
    "data": "<json data>",
    "output": "h264_baseline_360p_600.mp4",
    "active": true
  },
  "status": 200
}
```

---

####  Update Preset
```
POST /api/preset/:preset_id
```

##### Parameters
```
Content-Type: application/json
```

```json
{
    "name": "preset name",
    "description": "preset description",
    "data": "<json string of ffmpeg preset data>",
    "active": true
}
```

##### Response
```
Content-Type: application/json
```

```json
{
  "id": 1,
  "name": "preset name",
  "description": "preset description",
  "data": "<json string of ffmpeg preset data>",
  "active": true
}
```
