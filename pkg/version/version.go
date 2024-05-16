package version

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
	if oldV.Major != newV.Major {
		return oldV.Major < newV.Major
	}
	if oldV.Minor != newV.Minor {
		return oldV.Minor < newV.Minor
	}

	return oldV.Patch < newV.Patch
}
