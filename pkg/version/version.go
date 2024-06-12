package version

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

// Parse takes a version string in the format "Major.Minor.Patch" and populates the [*Version] struct with the parsed values.
// It returns an error if the input string is not in the correct format or if the numeric values cannot be parsed.
func Parse(s string) (*Version, error) {
	// split input string
	parts := strings.Split(s, ".")

	// check for consistency
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid version string: %s", s)
	}

	// parse version parts
	major, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid major version")
	}
	minor, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", s)
	}
	patch, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", s)
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

// String return a string representation of the [*Version] object in the format "Major.Minor.Patch".
func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Compare compares the reference [*Version] object against another [*Version] object.
// It returns -1 if the reference version is older, 0 if they are equal, and 1 if the reference version is newer.
func (v1 *Version) Compare(v2 *Version) int8 {
	// compare major version
	if v1.Major != v2.Major {
		if v1.Major < v2.Major {
			return -1
		}
		return 1
	}

	// compare minor version
	if v1.Minor != v2.Minor {
		if v1.Minor < v2.Minor {
			return -1
		}
		return 1
	}

	// compare patch
	if v1.Patch != v2.Patch {
		if v1.Patch < v2.Patch {
			return -1
		}
		return 1
	}

	// versions match
	return 0
}

// OlderThan returns true only if the reference version is older than the provided version, and false otherwise.
func (v1 *Version) OlderThan(v2 *Version) bool {
	return v1.Compare(v2) == -1
}

// Equals returns true if the reference version is equal to the provided version, and false otherwise.
func (v1 *Version) Equals(v2 *Version) bool {
	return v1.Compare(v2) == 0
}

// NewerThan returns true if the reference version is newer than the provided version, and false otherwise.
func (v1 *Version) NewerThan(v2 *Version) bool {
	return v1.Compare(v2) == 1
}
