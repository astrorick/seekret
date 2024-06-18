package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/astrorick/seekret/pkg/version"
)

type Database struct {
	SQL *sql.DB
}

func Open(databaseType string, databaseConnStr string, appVersion *version.Version) (*Database, error) {
	// when using an 'sqlite3' database, the database file must be created if it does not exist
	if databaseType == "sqlite3" {
		if _, err := os.Stat(databaseConnStr); errors.Is(err, os.ErrNotExist) {
			if _, err := os.Create(databaseConnStr); err != nil {
				return nil, err
			}
		}
	}

	// open connection to database
	sqlDB, err := sql.Open(databaseType, databaseConnStr)
	if err != nil {
		return nil, err
	}

	// new db object
	db := &Database{
		SQL: sqlDB,
	}

	// assess database status
	var statsRow *sql.Row
	var usersRow *sql.Row
	switch databaseType {
	case "sqlite3":
		statsRow = sqlDB.QueryRow("SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'stats'")
		usersRow = sqlDB.QueryRow("SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'users'")
	default:
		return nil, errors.New("unsupported database type")
	}

	// check the 'stats' table
	if err := statsRow.Scan(); errors.Is(err, sql.ErrNoRows) {
		// generate empty 'stats' table
		if _, err := sqlDB.Exec("CREATE TABLE stats (id INTEGER PRIMARY KEY, version TEXT NOT NULL)"); err != nil {
			return nil, err
		}

		// write default values
		if _, err := sqlDB.Exec("INSERT INTO stats (version) VALUES (?)", appVersion.String()); err != nil {
			return nil, err
		}
	} else {
		// read first row of 'stats' table
		stat, err := db.GetStat()
		if err != nil {
			return nil, err
		}

		// check versions mismatch
		if appVersion.OlderThan(stat.Version) {
			// database was created with a newer server version
			return nil, fmt.Errorf("outdated server version (%s) for the provided database (%s)", appVersion.String(), stat.Version.String())
		}
		if appVersion.NewerThan(stat.Version) {
			/*
				TODO
				We should proceed with migrations at this point. For now, if versions do not match the app simply stops.
			*/

			return nil, fmt.Errorf("migrations not implemented yet (server version: %s, database version: %s)", appVersion.String(), stat.Version.String())
		}
	}

	// check the 'users' table
	if err := usersRow.Scan(); errors.Is(err, sql.ErrNoRows) {
		// generate empty 'users' table
		if _, err := sqlDB.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT NOT NULL, salt TEXT NOT NULL, verifier TEXT NOT NULL)"); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func (db *Database) Close() error {
	if err := db.SQL.Close(); err != nil {
		return err
	}

	return nil
}
