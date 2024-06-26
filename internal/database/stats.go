package database

import "github.com/astrorick/seekret/pkg/version"

type Stat struct {
	ID      uint64           `db:"id"`
	Version *version.Version `db:"version"`
}

func (db *Database) GetStat() (*Stat, error) {
	var (
		stat          Stat
		versionString string
	)

	// select the first row of the 'stats' table
	if err := db.SQL.QueryRow("SELECT * FROM stats").Scan(&stat.ID, &versionString); err != nil {
		return nil, err
	}

	// parse version string
	version, err := version.New(versionString)
	if err != nil {
		return nil, err
	}
	stat.Version = version

	// return
	return &stat, nil
}
