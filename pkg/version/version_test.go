package version_test

import (
	"testing"

	"github.com/astrorick/seekret/pkg/version"
)

// TestParse calls version.Parse() with different values and checks for valid returns
func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		versionString   string
		expectedVersion *version.Version
		errorIsExpected bool
	}{
		// successfull tests
		{
			name:          "TestParse_0.0.0",
			versionString: "0.0.0",
			expectedVersion: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			errorIsExpected: false,
		},
		{
			name:          "TestParse_1.0.0",
			versionString: "1.0.0",
			expectedVersion: &version.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			errorIsExpected: false,
		},
		{
			name:          "TestParse_0.2.0",
			versionString: "0.2.0",
			expectedVersion: &version.Version{
				Major: 0,
				Minor: 2,
				Patch: 0,
			},
			errorIsExpected: false,
		},
		{
			name:          "TestParse_0.0.3",
			versionString: "0.0.3",
			expectedVersion: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 3,
			},
			errorIsExpected: false,
		},
		{
			name:          "TestParse_1.2.3",
			versionString: "1.2.3",
			expectedVersion: &version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			errorIsExpected: false,
		},

		// characters in input string
		{
			name:            "TestParse_a.2.3",
			versionString:   "a.2.3",
			expectedVersion: nil,
			errorIsExpected: true,
		},
		{
			name:            "TestParse_1.b.3",
			versionString:   "1.b.3",
			expectedVersion: nil,
			errorIsExpected: true,
		},
		{
			name:            "TestParse_1.2.c",
			versionString:   "1.2.c",
			expectedVersion: nil,
			errorIsExpected: true,
		},
		{
			name:            "TestParse_a.b.c",
			versionString:   "a.b.c",
			expectedVersion: nil,
			errorIsExpected: true,
		},

		// negative numbers
		{
			name:            "TestParse_-1.0.0",
			versionString:   "-1.0.0",
			expectedVersion: nil,
			errorIsExpected: true,
		},
		{
			name:            "TestParse_0.-1.0",
			versionString:   "0.-1.0",
			expectedVersion: nil,
			errorIsExpected: true,
		},
		{
			name:            "TestParse_0.0.-1",
			versionString:   "0.0.-1",
			expectedVersion: nil,
			errorIsExpected: true,
		},

		// garbage strings
		{
			name:            "TestParse_...",
			versionString:   "...",
			expectedVersion: nil,
			errorIsExpected: true,
		},
		{
			name:            "TestParse_1a.2b.3c",
			versionString:   "1a.2b.3c",
			expectedVersion: nil,
			errorIsExpected: true,
		},
		{
			name:            "TestParse_hello",
			versionString:   "hello",
			expectedVersion: nil,
			errorIsExpected: true,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			res, err := version.Parse(test.versionString)

			if err != nil {
				if (err != nil) != test.errorIsExpected {
					t.Errorf("expected: %v, result: %v", test.errorIsExpected, err != nil)
				}
				return
			}

			if *res != *test.expectedVersion {
				t.Errorf("expected: %v, result: %v", test.expectedVersion, res)
			}
		})
	}
}

// TestString calls Version.String() with different values and checks for valid returns
func TestString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		version        *version.Version
		expectedString string
	}{
		{
			name: "TestString_0.0.0",
			version: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			expectedString: "0.0.0",
		},
		{
			name: "TestString_1.0.0",
			version: &version.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			expectedString: "1.0.0",
		},
		{
			name: "TestString_0.2.0",
			version: &version.Version{
				Major: 0,
				Minor: 2,
				Patch: 0,
			},
			expectedString: "0.2.0",
		},
		{
			name: "TestString_0.0.3",
			version: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 3,
			},
			expectedString: "0.0.3",
		},
		{
			name: "TestString_1.2.3",
			version: &version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			expectedString: "1.2.3",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if res := test.version.String(); res != test.expectedString {
				t.Errorf("expected: %v, result: %v", test.expectedString, res)
			}
		})
	}
}

// TODO docs
func TestCompare(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		referenceVersion *version.Version
		compareVersion   *version.Version
		expectedInt      int8
	}{
		// tests with referenceVersion < compareVersion
		{
			name: "TestCompare_0.0.0_0.0.1",
			referenceVersion: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			compareVersion: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
			},
			expectedInt: -1,
		},
		{
			name: "TestCompare_1.5.7_1.6.6",
			referenceVersion: &version.Version{
				Major: 1,
				Minor: 5,
				Patch: 7,
			},
			compareVersion: &version.Version{
				Major: 1,
				Minor: 6,
				Patch: 6,
			},
			expectedInt: -1,
		},
		{
			name: "TestCompare_2.3.4_4.3.2",
			referenceVersion: &version.Version{
				Major: 2,
				Minor: 3,
				Patch: 4,
			},
			compareVersion: &version.Version{
				Major: 4,
				Minor: 3,
				Patch: 2,
			},
			expectedInt: -1,
		},

		// tests with referenceVersion = compareVersion
		{
			name: "TestCompare_1.2.3",
			referenceVersion: &version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			compareVersion: &version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			expectedInt: 0,
		},
		{
			name: "TestCompare_1.1.1",
			referenceVersion: &version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			compareVersion: &version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			expectedInt: 0,
		},
		{
			name: "TestCompare_1.11.111",
			referenceVersion: &version.Version{
				Major: 1,
				Minor: 11,
				Patch: 111,
			},
			compareVersion: &version.Version{
				Major: 1,
				Minor: 11,
				Patch: 111,
			},
			expectedInt: 0,
		},

		// tests with referenceVersion > compareVersion
		{
			name: "TestCompare_0.0.1_0.0.0",
			referenceVersion: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
			},
			compareVersion: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			expectedInt: 1,
		},
		{
			name: "TestCompare_1.2.3_1.1.4",
			referenceVersion: &version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			compareVersion: &version.Version{
				Major: 1,
				Minor: 1,
				Patch: 4,
			},
			expectedInt: 1,
		},
		{
			name: "TestCompare_3.1.1_2.5.5",
			referenceVersion: &version.Version{
				Major: 3,
				Minor: 1,
				Patch: 1,
			},
			compareVersion: &version.Version{
				Major: 2,
				Minor: 5,
				Patch: 5,
			},
			expectedInt: 1,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if res := test.referenceVersion.Compare(test.compareVersion); res != test.expectedInt {
				t.Errorf("expected: %v, result: %v", test.expectedInt, res)
			}
		})
	}
}

// TODO docs
func TestOlderThan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		referenceVersion *version.Version
		compareVersion   *version.Version
		expectedBool     bool
	}{
		{
			// tests with referenceVersion < compareVersion
			name: "TestOlderThan_1.5.5_2.5.2",
			referenceVersion: &version.Version{
				Major: 1,
				Minor: 5,
				Patch: 5,
			},
			compareVersion: &version.Version{
				Major: 2,
				Minor: 5,
				Patch: 2,
			},
			expectedBool: true,
		},

		// tests with referenceVersion = compareVersion
		{
			name: "TestOlderThan_8.9.4_8.9.4",
			referenceVersion: &version.Version{
				Major: 8,
				Minor: 9,
				Patch: 4,
			},
			compareVersion: &version.Version{
				Major: 8,
				Minor: 9,
				Patch: 4,
			},
			expectedBool: false,
		},

		// tests with referenceVersion > compareVersion
		{
			name: "TestOlderThan_9.5.1_8.4.3",
			referenceVersion: &version.Version{
				Major: 9,
				Minor: 5,
				Patch: 1,
			},
			compareVersion: &version.Version{
				Major: 8,
				Minor: 4,
				Patch: 3,
			},
			expectedBool: false,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if res := test.referenceVersion.OlderThan(test.compareVersion); res != test.expectedBool {
				t.Errorf("expected: %v, result: %v", test.expectedBool, res)
			}
		})
	}
}

// TODO docs
func TestEquals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		referenceVersion *version.Version
		compareVersion   *version.Version
		expectedBool     bool
	}{
		{
			// tests with referenceVersion < compareVersion
			name: "TestEquals_1.5.5_2.5.2",
			referenceVersion: &version.Version{
				Major: 1,
				Minor: 5,
				Patch: 5,
			},
			compareVersion: &version.Version{
				Major: 2,
				Minor: 5,
				Patch: 2,
			},
			expectedBool: false,
		},

		// tests with referenceVersion = compareVersion
		{
			name: "TestEquals_8.9.4_8.9.4",
			referenceVersion: &version.Version{
				Major: 8,
				Minor: 9,
				Patch: 4,
			},
			compareVersion: &version.Version{
				Major: 8,
				Minor: 9,
				Patch: 4,
			},
			expectedBool: true,
		},

		// tests with referenceVersion > compareVersion
		{
			name: "TestEquals_9.5.1_8.4.3",
			referenceVersion: &version.Version{
				Major: 9,
				Minor: 5,
				Patch: 1,
			},
			compareVersion: &version.Version{
				Major: 8,
				Minor: 4,
				Patch: 3,
			},
			expectedBool: false,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if res := test.referenceVersion.Equals(test.compareVersion); res != test.expectedBool {
				t.Errorf("expected: %v, result: %v", test.expectedBool, res)
			}
		})
	}
}

// TODO docs
func TestNewerThan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		referenceVersion *version.Version
		compareVersion   *version.Version
		expectedBool     bool
	}{
		{
			// tests with referenceVersion < compareVersion
			name: "TestEquals_1.5.5_2.5.2",
			referenceVersion: &version.Version{
				Major: 1,
				Minor: 5,
				Patch: 5,
			},
			compareVersion: &version.Version{
				Major: 2,
				Minor: 5,
				Patch: 2,
			},
			expectedBool: false,
		},

		// tests with referenceVersion = compareVersion
		{
			name: "TestEquals_8.9.4_8.9.4",
			referenceVersion: &version.Version{
				Major: 8,
				Minor: 9,
				Patch: 4,
			},
			compareVersion: &version.Version{
				Major: 8,
				Minor: 9,
				Patch: 4,
			},
			expectedBool: false,
		},

		// tests with referenceVersion > compareVersion
		{
			name: "TestEquals_9.5.1_8.4.3",
			referenceVersion: &version.Version{
				Major: 9,
				Minor: 5,
				Patch: 1,
			},
			compareVersion: &version.Version{
				Major: 8,
				Minor: 4,
				Patch: 3,
			},
			expectedBool: true,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if res := test.referenceVersion.NewerThan(test.compareVersion); res != test.expectedBool {
				t.Errorf("expected: %v, result: %v", test.expectedBool, res)
			}
		})
	}
}
