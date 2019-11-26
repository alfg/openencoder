package data

// Data represents the available database tables.
type Data struct {
	Presets  Presets
	Settings Settings
	Jobs     Jobs
	Users    Users
}

// New creates a new database instance.
func New() *Data {
	return &Data{
		Presets:  &PresetsOp{},
		Settings: &SettingsOp{},
		Jobs:     &JobsOp{},
		Users:    &UsersOp{},
	}
}
