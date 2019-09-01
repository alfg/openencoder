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

func CreateSettingByUserID() {

	return
}

func UpdateSettingByUserID() {

}

func UpdateSettingsByUserID() {

}
