package seekret

import "fmt"

type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (oldV *Version) IsOlderThan(newV *Version) bool {
	// check major
	if oldV.Major != newV.Major {
		return oldV.Major < newV.Major
	}

	// check minor
	if oldV.Minor != newV.Minor {
		return oldV.Minor < newV.Minor
	}

	// check patch
	return oldV.Patch < newV.Patch
}
