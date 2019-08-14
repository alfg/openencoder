package data

import (
	"fmt"

	"github.com/alfg/openencoder/api/types"
)

// CreateUser creates a user.
func CreateUser(user types.User) *types.User {
	const query = "INSERT INTO user (email) VALUES (:email)"

	db, _ := ConnectDB()
	tx := db.MustBegin()
	result, err := tx.NamedExec(query, &user)
	if err != nil {
		fmt.Println("Error", err)
	}
	tx.Commit()

	fmt.Println("transaction done")

	lastID, _ := result.LastInsertId()
	user.ID = lastID

	return &user
}
