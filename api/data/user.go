package data

import (
	"fmt"

	"github.com/alfg/openencoder/api/types"
)

// CreateUser creates a user.
func CreateUser(user types.User) (*types.User, error) {
	const query = `
	  INSERT INTO
	    users (username,password,role)
	  VALUES (:username,:password,:role)
	  RETURNING id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		fmt.Println("Error", err)
	}

	var id int64 // Returned ID.
	err = stmt.QueryRowx(&user).Scan(&id)
	if err != nil {
		fmt.Println("Error", err.Error())
		return nil, err
	}
	tx.Commit()
	user.ID = id

	db.Close()
	return &user, nil
}

// GetUserByUsername Gets a user by username.
func GetUserByUsername(username string) (*types.User, error) {
	const query = `
      SELECT
        users.*
      FROM users
      WHERE users.username = $1`

	db, _ := ConnectDB()
	user := types.User{}
	err := db.Get(&user, query, username)
	if err != nil {
		fmt.Println(err)
		return &user, err
	}
	db.Close()
	return &user, nil
}

// GetUserID Gets a user ID by username.
func GetUserID(username string) int64 {
	const query = "SELECT id FROM users WHERE username = $1"

	var id int64

	db, _ := ConnectDB()
	err := db.QueryRow(query, username).Scan(&id)
	if err != nil {
		fmt.Println(err)
	}
	return id
}
