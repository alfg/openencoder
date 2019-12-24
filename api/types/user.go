package types

// User contains user models.
type User struct {
	ID                 int64  `db:"id" json:"id,omitempty"`
	Username           string `db:"username" json:"username" valid:"email,required"`
	Password           string `db:"password" json:"password" valid:"password,required"`
	Role               string `db:"role" json:"role" valid:"required"`
	ForcePasswordReset bool   `db:"force_password_reset" json:"force_password_reset"`
}
