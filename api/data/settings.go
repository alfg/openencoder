package data

import (
	"fmt"

	"github.com/alfg/openencoder/api/types"
)

// GetSettingsByUserID Gets settings by a UserID.
func GetSettingsByUserID(id int64) []types.Setting {
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

	for k, v := range s {
		fmt.Println("run db update: ", k, v)
	}
	return nil
}
