package types

// Settings types.
const (
	AWSAccessKey      = "AWS_ACCESS_KEY"
	AWSSecretKey      = "AWS_SECRET_KEY"
	DigitalOceanToken = "DIGITAL_OCEAN_TOKEN"
)

// SettingsTypes list of all settings available.
var SettingsTypes = []string{
	AWSAccessKey,
	AWSSecretKey,
	DigitalOceanToken,
}

// Setting defines a setting for a user.
type Setting struct {
	ID     int64  `db:"id" json:"-"`
	UserID string `db:"user_id" json:"-"`

	SettingsOptionID int64 `db:"settings_option_id" json:"-"`
	SettingsOption   `db:"settings_option"`

	Value string `db:"value" json:"value"`
}

// SettingsOption defines a setting option.
type SettingsOption struct {
	ID          int64  `db:"id" json:"-"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}
