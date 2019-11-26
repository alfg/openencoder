package types

// Preset contains user models.
type Preset struct {
	ID          int64      `db:"id" json:"id,omitempty"`
	Name        string     `db:"name" json:"name" valid:"required"`
	Description NullString `db:"description" json:"description"`
	Options     NullString `db:"options" json:"options" valid:"required"`
	Active      bool       `db:"active" json:"active"`
}
