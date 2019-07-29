package data

import (
	"fmt"

	"github.com/alfg/openencoder/api/types"
)

// GetJobs Gets all jobs.
func GetJobs(offset, count int) *[]types.Job {
	const query = `SELECT * FROM jobs ORDER BY id DESC
    LIMIT $1 OFFSET $2`

	db, _ := ConnectDB()
	jobs := []types.Job{}
	err := db.Select(&jobs, query, count, offset)
	if err != nil {
		fmt.Println(err)
	}
	db.Close()
	return &jobs
}

// GetJobByID Gets a job by ID.
func GetJobByID(id int) (*types.Job, error) {
	const query = `
      SELECT
        jobs.*,
        encode_data.id "encode_data.id",
        encode_data.data "encode_data.data",
        encode_data.progress "encode_data.progress"
      FROM jobs
      LEFT JOIN encode_data ON jobs.id = encode_data.job_id
      WHERE jobs.id = $1`

	db, _ := ConnectDB()
	job := types.Job{}
	err := db.Get(&job, query, id)
	if err != nil {
		fmt.Println(err)
		return &job, err
	}
	db.Close()
	return &job, nil
}

// GetJobByGUID Gets a job by GUID.
func GetJobByGUID(id string) (*types.Job, error) {
	const query = `
      SELECT
        jobs.*,
        encode_data.id "encode_data.id",
        encode_data.data "encode_data.data",
        encode_data.progress "encode_data.progress"
      FROM jobs
      LEFT JOIN encode_data ON jobs.id = encode_data.job_id
      WHERE jobs.guid = $1`

	db, _ := ConnectDB()
	job := types.Job{}
	err := db.Get(&job, query, id)
	if err != nil {
		fmt.Println(err)
		return &job, err
	}
	db.Close()
	return &job, nil
}

// GetJobsCount Gets a count of all jobs.
func GetJobsCount() int {
	var count int
	const query = `SELECT COUNT(*) FROM jobs`

	db, _ := ConnectDB()
	err := db.Get(&count, query)
	if err != nil {
		fmt.Println(err)
	}
	db.Close()
	return count
}

// Stats struct for displaying status and count of a job.
type Stats struct {
	Status string `db:"status" json:"status"`
	Count  int    `db:"count" json:"count"`
}

// GetJobsStats Gets a count of each status.
func GetJobsStats() (*[]Stats, error) {
	const query = `SELECT status, count(status) FROM jobs GROUP BY status, status;`

	s := []Stats{}
	db, _ := ConnectDB()
	err := db.Select(&s, query)
	if err != nil {
		fmt.Println(err)
		return &s, err
	}
	db.Close()

	// Set all statuses.
	var resp []Stats
	for _, v := range types.JobStatuses {
		r := Stats{}
		for _, j := range s {
			if j.Status == v {
				r.Status = j.Status
				r.Count = j.Count
			} else {
				r.Status = v
			}
		}
		resp = append(resp, r)
	}
	return &resp, nil
}

// CreateJob creates a job in database.
func CreateJob(job types.Job) *types.Job {
	const query = `
      INSERT INTO
        jobs (guid,profile,status)
      VALUES (:guid,:profile,:status)
      RETURNING id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		fmt.Println("Error", err.Error())
	}

	var id int64 // Returned ID.
	err = stmt.QueryRowx(&job).Scan(&id)
	if err != nil {
		fmt.Println("Error", err.Error())
	}
	tx.Commit()

	// Set to Job type response.
	job.ID = id

	db.Close()
	return &job
}

// CreateEncodeData creates encode_data in database.
func CreateEncodeData(ed types.EncodeData) *types.EncodeData {
	const query = `
      INSERT INTO
        encode_data (data,progress,job_id)
      VALUES (:data,:progress,:job_id)
      RETURNING id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		fmt.Println("Error", err.Error())
	}

	var id int64 // Returned ID.
	err = stmt.QueryRowx(&ed).Scan(&id)
	if err != nil {
		fmt.Println("Error", err.Error())
	}
	tx.Commit()

	// Set to Job type response.
	ed.EncodeDataID = id

	db.Close()
	return &ed
}

// UpdateEncodeDataByID Update encode_data by ID.
func UpdateEncodeDataByID(id int64, jsonString string) error {
	const query = `UPDATE encode_data SET data = $1 WHERE id = $2`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.Exec(query, jsonString, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	tx.Commit()

	db.Close()
	return nil
}

// UpdateEncodeProgressByID Update progress by ID.
func UpdateEncodeProgressByID(id int64, progress float64) error {
	const query = `UPDATE encode_data SET progress = $1 WHERE id = $2`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.Exec(query, progress, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	tx.Commit()

	db.Close()
	return nil
}

// UpdateJobByID Update job by ID.
func UpdateJobByID(id int, job types.Job) *types.Job {
	const query = `UPDATE jobs SET status = :status WHERE id = :id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.NamedExec(query, &job)
	if err != nil {
		fmt.Println(err)
	}
	tx.Commit()

	db.Close()
	return &job
}

// UpdateJobStatus Update job status by ID.
func UpdateJobStatus(guid string, status string) error {
	const query = `UPDATE jobs SET status = $1 WHERE guid = $2`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.Exec(query, status, guid)
	if err != nil {
		fmt.Println(err)
		return err
	}
	tx.Commit()

	db.Close()
	return nil
}
