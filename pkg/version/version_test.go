package version_test

import (
	"testing"

	"github.com/astrorick/seekret/pkg/version"
)

// TestVersionString calls Version.String() with different values and checks for valid return values
func TestVersionString(t *testing.T) {
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

func TestVersionComparison(t *testing.T) {
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

			if got := tc.oldVersion.IsOlderThan(tc.newVersion); got != tc.want {
				t.Errorf("Got %v, wanted %v.", got, tc.want)
			}
		})
	}
}
