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
	return &jobs
}

// GetJobByID Gets a job by ID.
func GetJobByID(id int) (*types.Job, error) {
	const query = `SELECT * FROM jobs WHERE id = $1`

	db, _ := ConnectDB()
	job := types.Job{}
	err := db.Get(&job, query, id)
	if err != nil {
		fmt.Println(err)
		return &job, err
	}
	return &job, nil
}

// GetJobsCount Gets a count of all jobs.
func GetJobsCount() int {
	var count int
	const query = `SELECT COUNT(*) FROM jobs`

	db, _ := ConnectDB()
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	return count
}

// CreateJob creates a job in database.
func CreateJob(job types.Job) *types.Job {
	const query = `INSERT INTO jobs (guid,profile,status) VALUES (:guid,:profile,:status)`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	result, err := tx.NamedExec(query, &job)

	if err != nil {
		fmt.Println("Error", err)
	}
	tx.Commit()

	fmt.Println("transaction done")

	lastID, _ := result.LastInsertId()
	job.ID = lastID

	return &job
}

// UpdateJobByID Update job by ID.
func UpdateJobByID(id int, job types.Job) *types.Job {
	const query = `UPDATE jobs
		SET status = :status
		WHERE id = :id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.NamedExec(query, &job)
	if err != nil {
		fmt.Println(err)
	}
	tx.Commit()

	return &job
}
