package database

import "github.com/astrorick/seekret/pkg/version"

type Stat struct {
	ID      uint64           `db:"id"`
	Version *version.Version `db:"version"`
}
