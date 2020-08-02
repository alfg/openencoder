package data

import (
	"encoding/hex"
	"fmt"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/helpers"
	"github.com/alfg/openencoder/api/types"
)

// Settings represents the Settings database operations.
type Settings interface {
	GetSetting(key string) types.Setting
	GetSettings() []types.Setting
	GetSettingsOptions() []types.SettingsOption
	CreateSetting(setting types.Setting) *types.Setting
	CreateOrUpdateSetting(key, value string)
	UpdateSettings(setting map[string]string) error
	UpdateSetting(setting types.Setting) *types.Setting
	SettingExists(optionID int64) bool
}

// SettingsOp represents the settings.
type SettingsOp struct {
	s *Settings
}

var _ Settings = &SettingsOp{}

// GetSetting Gets a setting.
func (s SettingsOp) GetSetting(key string) types.Setting {
	const query = `
	  SELECT
        settings.*,
        settings_option.id "settings_option.id",
        settings_option.name "settings_option.name",
        settings_option.title "settings_option.title",
        settings_option.description "settings_option.description",
        settings_option.secure "settings_option.secure"
	  FROM settings
      JOIN settings_option ON settings.settings_option_id = settings_option.id
      WHERE settings_option.name = $1
      ORDER BY id DESC`

	db, _ := ConnectDB()
	setting := types.Setting{}
	err := db.Get(&setting, query, key)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	if setting.Secure {
		enc, _ := hex.DecodeString(setting.Value)
		plaintext, _ := helpers.Decrypt(enc, config.Keyseed())
		setting.Value = string(plaintext)
	}
	return setting
}

// GetSettings Gets settings.
func (s SettingsOp) GetSettings() []types.Setting {
	const query = `
	  SELECT
        settings.*,
        settings_option.id "settings_option.id",
        settings_option.name "settings_option.name",
        settings_option.title "settings_option.title",
        settings_option.description "settings_option.description",
        settings_option.secure "settings_option.secure"
	  FROM settings
      JOIN settings_option ON settings.settings_option_id = settings_option.id
      ORDER BY id DESC`

	db, _ := ConnectDB()
	settings := []types.Setting{}
	err := db.Select(&settings, query)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	for i := range settings {
		if settings[i].Secure {
			enc, _ := hex.DecodeString(settings[i].Value)
			plaintext, _ := helpers.Decrypt(enc, config.Keyseed())
			settings[i].Value = string(plaintext)
		}
	}
	return settings
}

// GetSettingsOptions Gets all available setting options.
func (s SettingsOp) GetSettingsOptions() []types.SettingsOption {
	const query = "SELECT * FROM settings_option ORDER BY id ASC"

	db, _ := ConnectDB()
	options := []types.SettingsOption{}
	err := db.Select(&options, query)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
	return options
}

// CreateSetting Creates a setting.
func (s SettingsOp) CreateSetting(setting types.Setting) *types.Setting {
	const query = `
      INSERT INTO
        settings (settings_option_id,value,encrypted)
      VALUES (:settings_option_id,:value,:encrypted)
      RETURNING id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		log.Fatal(err.Error())
	}

	var id int64 // Returned ID.
	err = stmt.QueryRowx(&setting).Scan(&id)
	if err != nil {
		log.Fatal(err.Error())
	}
	tx.Commit()

	// Set to Job type response.
	setting.SettingsOptionID = id

	db.Close()
	return &setting
}

// CreateOrUpdateSetting Runs an "upsert"-like transaction for a setting.
func (s SettingsOp) CreateOrUpdateSetting(key, value string) {
	availableSettings := s.GetSettingsOptions()
	k := getOptionKeyID(availableSettings, key)
	isSecure := isSecure(availableSettings, key)
	exists := s.SettingExists(k)

	se := types.Setting{
		SettingsOptionID: k,
		Value:            value,
		Encrypted:        false,
	}

	if isSecure {
		ciphertext, _ := helpers.Encrypt([]byte(value), config.Keyseed())
		se.Value = fmt.Sprintf("%x", ciphertext)
		se.Encrypted = true
	}

	if exists {
		s.UpdateSetting(se)
	} else {
		s.CreateSetting(se)
	}
}

// UpdateSettings Updates settings.
func (s SettingsOp) UpdateSettings(setting map[string]string) error {

	// Run insert or update for each setting.
	for k, v := range setting {
		s.CreateOrUpdateSetting(k, v)
	}
	return nil
}

// UpdateSetting updates an existing setting.
func (s SettingsOp) UpdateSetting(setting types.Setting) *types.Setting {
	const query = `
        UPDATE settings
        SET value = :value, encrypted = :encrypted
        WHERE settings_option_id = :settings_option_id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.NamedExec(query, &setting)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()

	db.Close()
	return &setting
}

// SettingExists Queries a setting exists.
func (s SettingsOp) SettingExists(optionID int64) bool {
	const query = `
        SELECT EXISTS
        (SELECT id
        FROM settings
        WHERE settings_option_id = $1)`

	var exists bool
	db, _ := ConnectDB()
	err := db.QueryRow(query, optionID).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
	return exists
}

func getOptionKeyID(s []types.SettingsOption, key string) int64 {
	for _, a := range s {
		if a.Name == key {
			return a.ID
		}
	}
	return -1
}

func isSecure(s []types.SettingsOption, key string) bool {
	for _, a := range s {
		if a.Secure && key == a.Name {
			return true
		}
	}
	return false
}
