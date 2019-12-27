package data

import (
	"github.com/alfg/openencoder/api/types"
)

// Users represents the Users database operations.
type Users interface {
	CreateUser(user types.User) (*types.User, error)
	GetUserByUsername(username string) (*types.User, error)
	GetUserID(username string) int64
	UpdateUserByID(id int, user *types.User) (*types.User, error)
	UpdateUserPasswordByID(id int64, user *types.User) (*types.User, error)
	GetUsers(offset, count int) *[]types.User
	GetUserByID(id int) (*types.User, error)
	GetUsersCount() int
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
		log.Fatal(err)
	}

	var id int64 // Returned ID.
	err = stmt.QueryRowx(&user).Scan(&id)
	if err != nil {
		log.Fatal(err.Error())
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
		log.Fatal(err)
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
		log.Fatal(err)
	}
	return id
}

// UpdateUserByID Update user by ID.
func (u UsersOp) UpdateUserByID(id int, user *types.User) (*types.User, error) {
	const query = `
        UPDATE users
        SET username = :username, password = :password, active = :active, role = :role
        WHERE id = :id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.NamedExec(query, &user)
	if err != nil {
		log.Fatal(err)
		return user, err
	}
	tx.Commit()

	db.Close()
	return user, nil
}

// UpdateUserPasswordByID Update user password by ID and reset force_password_reset.
func (u UsersOp) UpdateUserPasswordByID(id int64, user *types.User) (*types.User, error) {
	const query = `
        UPDATE users
        SET password = :password, force_password_reset = false
        WHERE id = :id`

	db, _ := ConnectDB()
	tx := db.MustBegin()
	_, err := tx.NamedExec(query, &user)
	if err != nil {
		log.Fatal(err)
		return user, err
	}
	tx.Commit()

	db.Close()
	return user, nil
}

// GetUsers gets a list of users with an offset and count.
func (u UsersOp) GetUsers(offset, count int) *[]types.User {
	const query = `
	  SELECT *
	  FROM users 
	  ORDER BY id DESC
      LIMIT $1 OFFSET $2`

	db, _ := ConnectDB()
	users := []types.User{}
	err := db.Select(&users, query, count, offset)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
	return &users
}

// GetUserByID Gets a user by ID.
func (u UsersOp) GetUserByID(id int) (*types.User, error) {
	const query = `
      SELECT *
      FROM users
      WHERE id = $1`

	db, _ := ConnectDB()
	user := types.User{}
	err := db.Get(&user, query, id)
	if err != nil {
		return &user, err
	}
	db.Close()
	return &user, nil
}

// GetUsersCount Gets a count of all users.
func (u UsersOp) GetUsersCount() int {
	var count int
	const query = `SELECT COUNT(*) FROM users`

	db, _ := ConnectDB()
	err := db.Get(&count, query)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
	return count
}
