package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

type Database struct {
	SQL *sql.DB
}

func Open(databaseType string, databaseConnStr string) (*Database, error) {
	// when using a 'sqlite3' database, the database file must be created if it does not exist
	if databaseType == "sqlite3" {
		if _, err := os.Stat(databaseConnStr); errors.Is(err, os.ErrNotExist) {
			if _, err := os.Create(databaseConnStr); err != nil {
				log.Fatal(err)
			}
		}
	}

	// open connection to database
	database, err := sql.Open(databaseType, databaseConnStr)
	if err != nil {
		return nil, err
	}

	// TODO: run consistency checks

	return &Database{
		SQL: database,
	}, nil
}

func (db *Database) Close() error {
	if err := db.SQL.Close(); err != nil {
		return err
	}

	return nil
}

func (db *Database) UserCount() (uint64, error) {
	var userCount uint64
	if err := db.SQL.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
		return 0, err
	}

	return userCount, nil
}

/*
This function should be executed as a part of the server initialization procedure.
It is meant to run preliminary consistency checks on the provided database in order to initialize missing tables and set some default values.
*/
/*func (srv *server.Server) RunPreliminaryChecks() error {
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

	// check the 'stats' table
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
		// read database version from 'stats' table
		var databaseVersionString string
		srv.Database.QueryRow("SELECT version FROM stats").Scan(&databaseVersionString)

		// parse to 'Version' object
		databaseVersion, err := version.New(databaseVersionString)
		if err != nil {
			return err
		}

		// check versions mismatch
		if srv.Version.OlderThan(databaseVersion) {
			// database was created with a newer server version
			return fmt.Errorf("outdated server version (%s) for the provided database (%s)", srv.Version.String(), databaseVersion.String())
		}
		if srv.Version.NewerThan(databaseVersion) {
			// TODO: proceed to migration
			//return nil

			fmt.Printf("Database updated from version %s to version %s\n", databaseVersion.String(), srv.Version.String())
		}
	}

	// check the 'users' table
	if err := usersRow.Scan(); errors.Is(err, sql.ErrNoRows) {
		// generate empty 'users' table
		if _, err := srv.Database.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT NOT NULL, salt TEXT NOT NULL, verifier TEXT NOT NULL)"); err != nil {
			return err
		}
	}

	return nil
}*/
