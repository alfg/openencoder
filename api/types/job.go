package types

import "time"

type Job struct {
	ID               string
	Profile          string
	Source           string
	Destination      string
	LocalSource      string
	LocalDestination string
	Delay            time.Duration
}
