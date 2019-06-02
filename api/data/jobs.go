package data

import (
	"fmt"

	"github.com/alfg/enc/api/types"
)

// GetJobs Gets all jobs.
func GetJobs() *[]types.Job {
	const query = `SELECT * FROM jobs ORDER BY id ASC`

	db, _ := ConnectDB()
	jobs := []types.Job{}
	err := db.Select(&jobs, query)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(&jobs)

	return &jobs
}

// CreateJob creates a job in database.
func CreateJob(job types.Job) *types.Job {
	const query = "INSERT INTO jobs (guid,profile) VALUES (:guid,:profile)"

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
