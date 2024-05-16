package sekret

import (
	"fmt"

	"github.com/astrorick/sekret/pkg/version"
)

type Sekret struct {
	appVersion *version.Version
}

func New() *Sekret {
	return &Sekret{
		appVersion: &version.Version{
			Major: 0,
			Minor: 1,
			Patch: 0,
		},
	}
}

func (sek *Sekret) Start() error {
	fmt.Printf("Sekret v%s by Astrorick\n\n", sek.appVersion.String())

	return nil
}
