package types

// User contains user models.
type User struct {
	ID       int64  `db:"id" json:"id,omitempty"`
	Email    string `db:"email" json:"email,omitempty" valid:"email,required"`
	Password string
}
