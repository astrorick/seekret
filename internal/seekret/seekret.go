package seekret

import (
	"fmt"

	"github.com/astrorick/seekret/pkg/version"
)

type Sekret struct {
	appVersion *version.Version
}

func New() *Sekret {
	// define server version
	return &Sekret{
		appVersion: &version.Version{
			Major: 0,
			Minor: 2,
			Patch: 0,
		},
	}
}

func (sek *Sekret) Start() error {
	// print welcome string
	fmt.Printf("Seekret v%s by Astrorick\n\n", sek.appVersion.String())

	// exit
	return nil
}
