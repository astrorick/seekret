package seekret

import (
	"database/sql"
	"errors"
)

/*
This function should be executed as a part of the server initialization procedure.
It is meant to run preliminary consistency checks on the provided database in order to initialize missing tables and set some default values.
*/
func (srv *Server) runPreliminaryChecks() error {
	// checks
	var statsRow *sql.Row
	var usersRow *sql.Row

	// assess database status
	switch srv.Config.DatabaseType {
	case "sqlite3":
		statsRow = srv.Database.QueryRow("SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'stats'")
		usersRow = srv.Database.QueryRow("SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'users'")
	default:
		return errors.New("unsupported database type")
	}

	// check if tables exist
	if err := statsRow.Scan(); errors.Is(err, sql.ErrNoRows) {
		// generate empty 'stats' table
		if _, err := srv.Database.Exec("CREATE TABLE stats (id INTEGER PRIMARY KEY, version TEXT NOT NULL)"); err != nil {
			return err
		}

		// write default values
		if _, err := srv.Database.Exec("INSERT INTO stats (version) VALUES (?)", srv.Version.String()); err != nil {
			return err
		}
	} else {
		// TODO: read database version and proceed to migration if necessary
	}

	if err := usersRow.Scan(); errors.Is(err, sql.ErrNoRows) {
		// generate empty 'users' table
		if _, err := srv.Database.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT NOT NULL, salt TEXT NOT NULL, verifier TEXT NOT NULL)"); err != nil {
			return err
		}
	}

	return nil
}
