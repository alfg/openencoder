package types

// Preset contains user models.
type Preset struct {
	ID          int64  `db:"id" json:"id,omitempty"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Data        string `db:"data" json:"data"`
	Output      string `db:"output" json:"output"`
	Active      *bool  `db:"active" json:"active"`
}
