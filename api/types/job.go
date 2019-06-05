package types

type Job struct {
	ID               int64  `db:"id" json:"id,omitempty"`
	GUID             string `db:"guid" json:"guid,omitempty"`
	Profile          string `db:"profile" json:"profile,omitempty"`
	Source           string `json:"source,omitempty"`
	Destination      string `json:"destination,omitempty"`
	LocalSource      string `json:"local_source,omitempty"`
	LocalDestination string `json:"local_destination,omitempty"`
	CreatedDate      string `db:"created_date" json:"created_date,omitempty"`
}
