package types

// User contains user models.
type User struct {
	ID                 int64  `db:"id" json:"id,omitempty"`
	Username           string `db:"username" json:"username" valid:"required"`
	Password           string `db:"password" json:"-" valid:"password,required"`
	Role               string `db:"role" json:"role" valid:"required"`
	ForcePasswordReset bool   `db:"force_password_reset" json:"-"`
	Active             bool   `db:"active" json:"active"`
}
