package data

import (
	"fmt"

	"github.com/alfg/openencoder/api/types"
)

// Users represents the Users database operations.
type Users interface {
	CreateUser(user types.User) (*types.User, error)
	GetUserByUsername(username string) (*types.User, error)
	GetUserID(username string) int64
	UpdateUserByID(id int64, user *types.User) (*types.User, error)
}

// UsersOp represents the users operations.
type UsersOp struct {
	u *Users
}

var _ Users = &UsersOp{}

// CreateUser creates a user.
func (u UsersOp) CreateUser(user types.User) (*types.User, error) {
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
func (u UsersOp) GetUserByUsername(username string) (*types.User, error) {
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
func (u UsersOp) GetUserID(username string) int64 {
	const query = "SELECT id FROM users WHERE username = $1"

	var id int64

	db, _ := ConnectDB()
	err := db.QueryRow(query, username).Scan(&id)
	if err != nil {
		fmt.Println(err)
	}
	return id
}

// UpdateUserByID Update user by ID.
func (u UsersOp) UpdateUserByID(id int64, user *types.User) (*types.User, error) {
	const query = `
        UPDATE users
        SET username = :username, password = :password
        WHERE id = :id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.NamedExec(query, &user)
	if err != nil {
		fmt.Println(err)
		return user, err
	}
	tx.Commit()

	db.Close()
	return user, nil
}
