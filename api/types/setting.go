package types

// Settings types.
const (
	AWSAccessKey            = "AWS_ACCESS_KEY"
	AWSSecretKey            = "AWS_SECRET_KEY"
	DigitalOceanAccessToken = "DIGITAL_OCEAN_ACCESS_TOKEN"
	SlackWebhook            = "SLACK_WEBHOOK"
)

// SettingsTypes list of all settings available.
var SettingsTypes = []string{
	AWSAccessKey,
	AWSSecretKey,
	DigitalOceanAccessToken,
	SlackWebhook,
}

// Setting defines a setting for a user.
type Setting struct {
	ID     int64 `db:"id" json:"-"`
	UserID int64 `db:"user_id" json:"-"`

	SettingsOptionID int64 `db:"settings_option_id" json:"-"`
	SettingsOption   `db:"settings_option"`

	Value     string `db:"value" json:"value"`
	Encrypted bool   `db:"encrypted" json:"encrypted"`
}

// SettingsOption defines a setting option.
type SettingsOption struct {
	ID          int64  `db:"id" json:"-"`
	Name        string `db:"name" json:"name"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Secure      bool   `db:"secure" json:"secure"`
}

type SettingsForm struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Secure      bool   `json:"secure"`
}
