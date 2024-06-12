package version_test

import (
	"fmt"
	"testing"

	"github.com/astrorick/seekret/pkg/version"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		versionString   string
		expectedVersion *version.Version
		expectedError   error
	}{
		{
			name:          "TestParse-0.0.0",
			versionString: "0.0.0",
			expectedVersion: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			expectedError: nil,
		},
		{
			name:          "TestParse-1.0.0",
			versionString: "1.0.0",
			expectedVersion: &version.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			expectedError: nil,
		},
		{
			name:          "TestParse-0.2.0",
			versionString: "0.2.0",
			expectedVersion: &version.Version{
				Major: 0,
				Minor: 2,
				Patch: 0,
			},
			expectedError: nil,
		},
		{
			name:          "TestParse-0.0.3",
			versionString: "0.0.3",
			expectedVersion: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 3,
			},
			expectedError: nil,
		},
		{
			name:          "TestParse-1.2.3",
			versionString: "1.2.3",
			expectedVersion: &version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			expectedError: nil,
		},
		{
			name:            "TestParse-1.2.3",
			versionString:   "a.2.3",
			expectedVersion: nil,
			expectedError:   fmt.Errorf("invalid major version"),
		},
	}

	for _, test := range tests {
		res, err := version.Parse(test.versionString)
		if err != test.expectedError {
			t.Errorf("expected: %v, result: %v", test.expectedError, err)
		}
		if *res != *test.expectedVersion {
			t.Errorf("expected: %v, result: %v", test.expectedVersion, res)
		}
	}
}

// TestString() calls Version.String() with different values and checks for valid return values
func TestString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ver  *version.Version
		want string
	}{
		{
			name: "TestString-0.0.0",
			ver: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			want: "0.0.0",
		},
		{
			name: "TestString-1.0.0",
			ver: &version.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			want: "1.0.0",
		},
		{
			name: "TestString-0.2.0",
			ver: &version.Version{
				Major: 0,
				Minor: 2,
				Patch: 0,
			},
			want: "0.2.0",
		},
		{
			name: "TestString-0.0.3",
			ver: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 3,
			},
			want: "0.0.3",
		},
		{
			name: "TestString-1.2.3",
			ver: &version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			want: "1.2.3",
		},
	}

	for _, test := range tests {
		if got := test.ver.String(); got != test.want {
			t.Errorf("Got %v, wanted %v", got, test.want)
		}
	}
}

func TestCompare(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		oldVersion *version.Version
		newVersion *version.Version
		want       bool
	}{
		{
			name: "Comparison111211Test",
			oldVersion: &version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			newVersion: &version.Version{
				Major: 2,
				Minor: 1,
				Patch: 1,
			},
			want: true,
		},
		{
			name: "Comparison111121Test",
			oldVersion: &version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			newVersion: &version.Version{
				Major: 1,
				Minor: 2,
				Patch: 1,
			},
			want: true,
		},
		{
			name: "Comparison111112Test",
			oldVersion: &version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			newVersion: &version.Version{
				Major: 1,
				Minor: 1,
				Patch: 2,
			},
			want: true,
		},
		{
			name: "Comparison111222Test",
			oldVersion: &version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			newVersion: &version.Version{
				Major: 2,
				Minor: 2,
				Patch: 2,
			},
			want: true,
		},
		{
			name: "Comparison111011Test",
			oldVersion: &version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			newVersion: &version.Version{
				Major: 0,
				Minor: 1,
				Patch: 1,
			},
			want: false,
		},
		{
			name: "Comparison111101Test",
			oldVersion: &version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			newVersion: &version.Version{
				Major: 1,
				Minor: 0,
				Patch: 1,
			},
			want: false,
		},
		{
			name: "Comparison111110Test",
			oldVersion: &version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			newVersion: &version.Version{
				Major: 1,
				Minor: 1,
				Patch: 0,
			},
			want: false,
		},
		{
			name: "Comparison111000Test",
			oldVersion: &version.Version{
				Major: 1,
				Minor: 1,
				Patch: 1,
			},
			newVersion: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			want: false,
		},
		{
			name: "Comparison333342Test",
			oldVersion: &version.Version{
				Major: 3,
				Minor: 3,
				Patch: 3,
			},
			newVersion: &version.Version{
				Major: 3,
				Minor: 4,
				Patch: 2,
			},
			want: true,
		},
		{
			name: "Comparison333324Test",
			oldVersion: &version.Version{
				Major: 3,
				Minor: 3,
				Patch: 3,
			},
			newVersion: &version.Version{
				Major: 3,
				Minor: 2,
				Patch: 4,
			},
			want: false,
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.oldVersion.OlderThan(tc.newVersion); got != tc.want {
				t.Errorf("Got %v, wanted %v.", got, tc.want)
			}
		})
	}
}

func TestOlderThan(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		version *version.Version
		want    string
	}{
		{
			name: "String000Test",
			version: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			want: "0.0.0",
		},
		{
			name: "String100Test",
			version: &version.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			want: "1.0.0",
		},
		{
			name: "String010Test",
			version: &version.Version{
				Major: 0,
				Minor: 1,
				Patch: 0,
			},
			want: "0.1.0",
		},
		{
			name: "String001Test",
			version: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
			},
			want: "0.0.1",
		},
		{
			name: "String123Test",
			version: &version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			want: "1.2.3",
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.version.String(); got != tc.want {
				t.Errorf("Got %v, wanted %v", got, tc.want)
			}
		})
	}
}

func TestEquals(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		version *version.Version
		want    string
	}{
		{
			name: "String000Test",
			version: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			want: "0.0.0",
		},
		{
			name: "String100Test",
			version: &version.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			want: "1.0.0",
		},
		{
			name: "String010Test",
			version: &version.Version{
				Major: 0,
				Minor: 1,
				Patch: 0,
			},
			want: "0.1.0",
		},
		{
			name: "String001Test",
			version: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
			},
			want: "0.0.1",
		},
		{
			name: "String123Test",
			version: &version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			want: "1.2.3",
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.version.String(); got != tc.want {
				t.Errorf("Got %v, wanted %v", got, tc.want)
			}
		})
	}
}

func TestNewerThan(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		version *version.Version
		want    string
	}{
		{
			name: "String000Test",
			version: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			want: "0.0.0",
		},
		{
			name: "String100Test",
			version: &version.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			want: "1.0.0",
		},
		{
			name: "String010Test",
			version: &version.Version{
				Major: 0,
				Minor: 1,
				Patch: 0,
			},
			want: "0.1.0",
		},
		{
			name: "String001Test",
			version: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
			},
			want: "0.0.1",
		},
		{
			name: "String123Test",
			version: &version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			want: "1.2.3",
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.version.String(); got != tc.want {
				t.Errorf("Got %v, wanted %v", got, tc.want)
			}
		})
	}
}
