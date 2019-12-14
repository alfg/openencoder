package types

// Settings types.
const (
	S3AccessKey             = "S3_ACCESS_KEY"
	S3SecretKey             = "S3_SECRET_KEY"
	DigitalOceanAccessToken = "DIGITAL_OCEAN_ACCESS_TOKEN"
	SlackWebhook            = "SLACK_WEBHOOK"
	S3InboundBucket         = "S3_INBOUND_BUCKET"
	S3InboundBucketRegion   = "S3_INBOUND_BUCKET_REGION"
	S3OutboundBucket        = "S3_OUTBOUND_BUCKET"
	S3OutboundBucketRegion  = "S3_OUTBOUND_BUCKET_REGION"
	S3Provider              = "S3_PROVIDER"

	DigitalOcean = "DIGITALOCEAN"
	AWS          = "AWS"
)

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
