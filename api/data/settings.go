package data

import (
	"encoding/hex"
	"fmt"

	"github.com/alfg/openencoder/api/config"
	"github.com/alfg/openencoder/api/helpers"
	"github.com/alfg/openencoder/api/types"
)

// GetSettingByUsername Gets a setting by a UserID.
func GetSetting(key string) types.Setting {
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
func GetSettingsByUsername(username string) []types.Setting {
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
func GetSettingsOptions() []types.SettingsOption {
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
func UpdateSettingsByUserID(id int64, s map[string]string) error {

	// Run insert or update for each setting.
	for k, v := range s {
		CreateOrUpdateUserSetting(id, k, v)
	}
	return nil
}

// CreateOrUpdateUserSetting Runs an "upsert"-like transaction for a user setting.
func CreateOrUpdateUserSetting(id int64, key, value string) {
	availableSettings := GetSettingsOptions()
	k := getOptionKeyID(availableSettings, key)
	isSecure := isSecure(availableSettings)
	exists := SettingExists(id, k)

	s := types.Setting{
		SettingsOptionID: k,
		Value:            value,
		UserID:           id,
		Encrypted:        false,
	}

	if isSecure {
		ciphertext, _ := helpers.Encrypt([]byte(value), config.Keyseed())
		s.Value = fmt.Sprintf("%x", ciphertext)
		s.Encrypted = true
	}

	if exists {
		UpdateSetting(s)
	} else {
		CreateSetting(s)
	}
}

// UpdateSetting updates an existing user setting.
func UpdateSetting(s types.Setting) *types.Setting {
	const query = `
        UPDATE settings
        SET value = :value, encrypted = :encrypted
        WHERE user_id = :user_id
        AND settings_option_id = :settings_option_id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.NamedExec(query, &s)
	if err != nil {
		fmt.Println(err)
	}
	tx.Commit()

	db.Close()
	return &s
}

// CreateSetting Creates a user setting.
func CreateSetting(s types.Setting) *types.Setting {
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
	err = stmt.QueryRowx(&s).Scan(&id)
	if err != nil {
		fmt.Println("Error", err.Error())
	}
	tx.Commit()

	// Set to Job type response.
	s.SettingsOptionID = id

	db.Close()
	return &s
}

// SettingExists Queries a user setting exists.
func SettingExists(userID, optionID int64) bool {
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

func isSecure(s []types.SettingsOption) bool {
	for _, a := range s {
		if a.Secure {
			return true
		}
	}
	return false
}
