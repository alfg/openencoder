package types

// Settings types.
const (
	StorageDriver = "STORAGE_DRIVER"

	S3AccessKey            = "S3_ACCESS_KEY"
	S3SecretKey            = "S3_SECRET_KEY"
	S3InboundBucket        = "S3_INBOUND_BUCKET"
	S3InboundBucketRegion  = "S3_INBOUND_BUCKET_REGION"
	S3OutboundBucket       = "S3_OUTBOUND_BUCKET"
	S3OutboundBucketRegion = "S3_OUTBOUND_BUCKET_REGION"
	S3Provider             = "S3_PROVIDER"
	S3Endpoint             = "S3_ENDPOINT"
	S3Streaming            = "S3_STREAMING"

	FTPAddr     = "FTP_ADDR"
	FTPUsername = "FTP_USERNAME"
	FTPPassword = "FTP_PASSWORD"

	DigitalOceanEnabled     = "DIGITAL_OCEAN_ENABLED"
	DigitalOceanAccessToken = "DIGITAL_OCEAN_ACCESS_TOKEN"
	DigitalOceanRegion      = "DIGITAL_OCEAN_REGION"
	DigitalOceanVPC         = "DIGITAL_OCEAN_VPC"
	SlackWebhook            = "SLACK_WEBHOOK"

	DigitalOceanSpaces = "DIGITALOCEANSPACES"
	AmazonAWS          = "AMAZONAWS"
	Custom             = "CUSTOM"
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

// SettingsForm defines the setting form options.
type SettingsForm struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Secure      bool   `json:"secure"`
}

// GetSetting gets a setting value by key from a slice of Setting.
func GetSetting(s string, settings []Setting) string {
	for _, v := range settings {
		if s == v.Name {
			return v.Value
		}
	}
	return ""
}
