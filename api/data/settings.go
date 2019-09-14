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
	GetSettingsByUsername(username string) []types.Setting
	GetSettingsOptions() []types.SettingsOption
	UpdateSettingsByUserID(id int64, setting map[string]string) error
	CreateOrUpdateUserSetting(id int64, key, value string)
	UpdateSetting(setting types.Setting) *types.Setting
	CreateSetting(setting types.Setting) *types.Setting
	SettingExists(userID, optionID int64) bool
}

// SettingsOp represents the user settings.
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
		fmt.Println(err)
	}
	db.Close()

	if setting.Secure {
		enc, _ := hex.DecodeString(setting.Value)
		plaintext, _ := helpers.Decrypt(enc, config.Keyseed())
		setting.Value = string(plaintext)
	}
	return setting
}

// GetSettingsByUsername Gets settings by a UserID.
func (s SettingsOp) GetSettingsByUsername(username string) []types.Setting {
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
		fmt.Println(err)
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
		fmt.Println(err)
	}
	db.Close()
	return options
}

// UpdateSettingsByUserID Updates settings by UserID.
func (s SettingsOp) UpdateSettingsByUserID(id int64, setting map[string]string) error {

	// Run insert or update for each setting.
	for k, v := range setting {
		s.CreateOrUpdateUserSetting(id, k, v)
	}
	return nil
}

// CreateOrUpdateUserSetting Runs an "upsert"-like transaction for a user setting.
func (s SettingsOp) CreateOrUpdateUserSetting(id int64, key, value string) {
	availableSettings := s.GetSettingsOptions()
	k := getOptionKeyID(availableSettings, key)
	isSecure := isSecure(availableSettings, key)
	exists := s.SettingExists(id, k)

	se := types.Setting{
		SettingsOptionID: k,
		Value:            value,
		UserID:           id,
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

// UpdateSetting updates an existing user setting.
func (s SettingsOp) UpdateSetting(setting types.Setting) *types.Setting {
	const query = `
        UPDATE settings
        SET value = :value, encrypted = :encrypted
        WHERE user_id = :user_id
        AND settings_option_id = :settings_option_id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.NamedExec(query, &setting)
	if err != nil {
		fmt.Println(err)
	}
	tx.Commit()

	db.Close()
	return &setting
}

// CreateSetting Creates a user setting.
func (s SettingsOp) CreateSetting(setting types.Setting) *types.Setting {
	const query = `
      INSERT INTO
        settings (settings_option_id,value,user_id,encrypted)
      VALUES (:settings_option_id,:value,:user_id,:encrypted)
      RETURNING id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		fmt.Println("Error", err.Error())
	}

	var id int64 // Returned ID.
	err = stmt.QueryRowx(&setting).Scan(&id)
	if err != nil {
		fmt.Println("Error", err.Error())
	}
	tx.Commit()

	// Set to Job type response.
	setting.SettingsOptionID = id

	db.Close()
	return &setting
}

// SettingExists Queries a user setting exists.
func (s SettingsOp) SettingExists(userID, optionID int64) bool {
	const query = `
        SELECT EXISTS
        (SELECT id
        FROM settings
        WHERE user_id = $1 AND settings_option_id = $2)`

	var exists bool
	db, _ := ConnectDB()
	err := db.QueryRow(query, userID, optionID).Scan(&exists)
	if err != nil {
		fmt.Println(err)
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
