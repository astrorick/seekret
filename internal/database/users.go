package database

type User struct {
	ID       uint64 `db:"id"`
	Username string `db:"username"`
	Salt     string `db:"salt"`
	Verifier string `db:"verifier"`
}

func (db *Database) UserCount() (uint64, error) {
	var userCount uint64
	if err := db.SQL.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
		return 0, err
	}

	return userCount, nil
}

func (db *Database) UserExists(username string) (bool, error) {
	var userExists bool
	if err := db.SQL.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&userExists); err != nil {
		return false, err
	}

	return userExists, nil
}

func (db *Database) CreateUser(username string, salt []byte, verifier []byte) error {
	if _, err := db.SQL.Exec("INSERT INTO users (username, salt, verifier) VALUES (?, ?, ?)", username, salt, verifier); err != nil {
		return err
	}

	return nil
}

func (db *Database) GetUser(username string) (*User, error) {
	var user User
	if err := db.SQL.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Salt, &user.Verifier); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) DeleteUser(username string) error {
	if _, err := db.SQL.Exec("DELETE FROM users WHERE username = ?", username); err != nil {
		return err
	}

	return nil
}
