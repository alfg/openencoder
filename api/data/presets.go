package data

import (
	"github.com/alfg/openencoder/api/types"
)

// Presets represents the Presets database operations.
type Presets interface {
	GetPresets(offset, count int) *[]types.Preset
	GetPresetByID(id int) (*types.Preset, error)
	GetPresetByName(name string) (*types.Preset, error)
	GetPresetsCount() int
	CreatePreset(user types.Preset) (*types.Preset, error)
	UpdatePresetByID(id int, preset types.Preset) *types.Preset
	UpdatePresetStatusByID(id int, active bool) error
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
		log.Fatal(err)
	}
	db.Close()
	return &presets
}

// GetPresetByID Gets a preset by ID.
func (p PresetsOp) GetPresetByID(id int) (*types.Preset, error) {
	const query = `
      SELECT *
      FROM presets
      WHERE id = $1`

	db, _ := ConnectDB()
	preset := types.Preset{}
	err := db.Get(&preset, query, id)
	if err != nil {
		log.Fatal(err)
		return &preset, err
	}
	db.Close()
	return &preset, nil
}

// GetPresetByName Gets a preset by name.
func (p PresetsOp) GetPresetByName(name string) (*types.Preset, error) {
	const query = `
      SELECT *
      FROM presets
      WHERE name = $1`

	db, _ := ConnectDB()
	preset := types.Preset{}
	err := db.Get(&preset, query, name)
	if err != nil {
		log.Fatal(err)
		return &preset, err
	}
	db.Close()
	return &preset, nil
}

// GetPresetsCount Gets a count of all presets.
func (p PresetsOp) GetPresetsCount() int {
	var count int
	const query = `SELECT COUNT(*) FROM presets`

	db, _ := ConnectDB()
	err := db.Get(&count, query)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
	return count
}

// CreatePreset creates a preset.
func (p PresetsOp) CreatePreset(preset types.Preset) (*types.Preset, error) {
	const query = `
	  INSERT INTO
	    presets (name,description,data,active,output)
	  VALUES (:name,:description,:data,:active,:output)
	  RETURNING id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		log.Fatal(err)
	}

	var id int64 // Returned ID.
	err = stmt.QueryRowx(&preset).Scan(&id)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	tx.Commit()
	preset.ID = id

	db.Close()
	return &preset, nil
}

// UpdatePresetByID Update job by ID.
func (p PresetsOp) UpdatePresetByID(id int, preset types.Preset) *types.Preset {
	const query = `
        UPDATE presets
        SET name = :name, description = :description, data = :data, active = :active
        WHERE id = :id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.NamedExec(query, &preset)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()

	db.Close()
	return &preset
}

// UpdatePresetStatusByID Update preset status by ID.
func (p PresetsOp) UpdatePresetStatusByID(id int, active bool) error {
	const query = `UPDATE presets SET active = $2 WHERE id = $1`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.Exec(query, active, id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	tx.Commit()

	db.Close()
	return nil
}
