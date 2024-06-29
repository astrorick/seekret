package database

import "fmt"

type NewUser struct {
	Username string `db:"username"`
	Salt     []byte `db:"salt"`
	Verifier []byte `db:"verifier"`
}

type User struct {
	ID       uint64 `db:"id"`
	Username string `db:"username"`
	Salt     []byte `db:"salt"`
	Verifier []byte `db:"verifier"`
}

// UserCount returns the number of users in the 'users' table.
// A non-nil error is returned if the query cannot be executed.
func (db *Database) UserCount() (uint64, error) {
	var userCount uint64
	if err := db.SQL.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
		return 0, err
	}

	return userCount, nil
}

// UserExists returns true if the 'users' table contains a row with the specified username, and false otherwise.
// A non-nil error is returned if the query cannot be executed.
func (db *Database) UserExists(username string) (bool, error) {
	var userExists bool
	if err := db.SQL.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&userExists); err != nil {
		return false, err
	}

	return userExists, nil
}

// CreareUser reates a new user in the 'users' table.
// A non-nil error is returned if query execution fails.
func (db *Database) CreateUser(newUser *NewUser) error {
	if _, err := db.SQL.Exec("INSERT INTO users (username, salt, verifier) VALUES (?, ?, ?)", newUser.Username, newUser.Salt, newUser.Verifier); err != nil {
		return err
	}

	return nil
}

// GetUser reads the 'users' table and returns the user with the specified username.
// A non-nil error is returned if the query cannot be executed or if no user is found.
func (db *Database) GetUser(username string) (*User, error) {
	var user User
	if err := db.SQL.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Salt, &user.Verifier); err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser removes a user from the 'users' table based on its username.
// A non-nil error is returned if the query cannot be executed or no user is found.
func (db *Database) DeleteUser(username string) error {
	result, err := db.SQL.Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("found no user with username '%s'", username)
	}
	return nil
}
