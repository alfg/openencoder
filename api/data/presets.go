package data

import (
	"fmt"

	"github.com/alfg/openencoder/api/types"
)

// Presets represents the Presets database operations.
type Presets interface {
	GetPresets(offset, count int) *[]types.Preset
	GetPresetsCount() int
	CreatePreset(user types.Preset) (*types.Preset, error)
}

// PresetsOp represents the users operations.
type PresetsOp struct {
	u *Presets
}

var _ Presets = &PresetsOp{}

// GetPresets Gets all presets.
func (p PresetsOp) GetPresets(offset, count int) *[]types.Preset {
	const query = `
	  SELECT * FROM presets
	  ORDER BY id DESC
      LIMIT $1 OFFSET $2`

	db, _ := ConnectDB()
	presets := []types.Preset{}
	err := db.Select(&presets, query, count, offset)
	if err != nil {
		fmt.Println(err)
	}
	db.Close()
	return &presets
}

// GetPresetsCount Gets a count of all presets.
func (p PresetsOp) GetPresetsCount() int {
	var count int
	const query = `SELECT COUNT(*) FROM presets`

	db, _ := ConnectDB()
	err := db.Get(&count, query)
	if err != nil {
		fmt.Println(err)
	}
	db.Close()
	return count
}

// CreatePreset creates a preset.
func (p PresetsOp) CreatePreset(preset types.Preset) (*types.Preset, error) {
	const query = `
	  INSERT INTO
	    presets (name)
	  VALUES (:name)
	  RETURNING id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		fmt.Println("Error", err)
	}

	var id int64 // Returned ID.
	err = stmt.QueryRowx(&preset).Scan(&id)
	if err != nil {
		fmt.Println("Error", err.Error())
		return nil, err
	}
	tx.Commit()
	preset.ID = id

	db.Close()
	return &preset, nil
}
