package seekret

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func StringToVersion(s string) (*Version, error) {
	// split input string
	parts := strings.Split(s, ".")

	// check for consistency
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid input string: %s", s)
	}

	// parse version parts
	major, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid input string: %s", s)
	}
	minor, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid input string: %s", s)
	}
	patch, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid input string: %s", s)
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
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
