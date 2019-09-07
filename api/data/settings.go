package data

import (
	"fmt"

	"github.com/alfg/openencoder/api/types"
)

// GetSettingsByUsername Gets settings by a UserID.
func GetSettingsByUsername(username string) []types.Setting {
	const query = `
	  SELECT
        settings.*,
        settings_option.id "settings_option.id",
        settings_option.name "settings_option.name",
        settings_option.title "settings_option.title",
        settings_option.description "settings_option.description"
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
	exists := SettingExists(id, k)

	s := types.Setting{
		SettingsOptionID: k,
		Value:            value,
		UserID:           id,
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
        SET value = :value
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
        settings (settings_option_id,value,user_id)
      VALUES (:settings_option_id,:value,:user_id)
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
