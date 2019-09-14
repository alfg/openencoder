package data

// Data represents the available database tables.
type Data struct {
	Settings Settings
	Jobs     Jobs
	Users    Users
}

// New creates a new database instance.
func New() *Data {
	return &Data{
		Settings: &SettingsOp{},
		Jobs:     &JobsOp{},
		Users:    &UsersOp{},
	}
}
