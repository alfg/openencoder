package types

import (
	"database/sql"
	"encoding/json"
)

// Job status types.
const (
	JobCreated     = "created"
	JobPending     = "pending"
	JobDownloading = "downloading"
	JobEncoding    = "encoding"
	JobUploading   = "uploading"
	JobCompleted   = "completed"
	JobError       = "error"
)

// JobStatuses All job status types.
var JobStatuses = []string{
	JobCreated,
	JobPending,
	JobDownloading,
	JobEncoding,
	JobUploading,
	JobCompleted,
	JobError,
}

// Job describes the job info.
type Job struct {
	ID               int64  `db:"id" json:"id,omitempty"`
	GUID             string `db:"guid" json:"guid,omitempty"`
	Profile          string `db:"profile" json:"profile,omitempty"`
	CreatedDate      string `db:"created_date" json:"created_date"`
	Status           string `db:"status" json:"status"`
	Source           string `json:"source,omitempty"`
	Destination      string `json:"destination,omitempty"`
	LocalSource      string `json:"local_source,omitempty"`
	LocalDestination string `json:"local_destination,omitempty"`
}

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}
