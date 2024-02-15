package src

import (
	"strings"
	"testing"
)

func Test_findLicense(t *testing.T) {
	testCases := []struct {
		name    string
		search  string
		want    string
		wantErr string
	}{
		{
			name:   "empty",
			search: "",
			want:   "",
		},
		{
			name:   "none",
			search: "none",
			want:   "",
		},
		{
			name:   "exact name",
			search: "mit",
			want:   "mit",
		},
		{
			name:   "exact match",
			search: "newbsd",
			want:   "bsd",
		},
		{
			name:   "different case",
			search: "GPLv3",
			want:   "gpl3",
		},
		{
			name:    "not found",
			search:  "does not exist",
			wantErr: "no license matching",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			want, err := findLicense(tC.search)
			if tC.wantErr != "" {
				if !strings.Contains(err.Error(), tC.wantErr) {
					t.Errorf("expected `%s` to contain `%s`", err.Error(), tC.wantErr)
				}
				return
			}
			if tC.want != want {
				t.Errorf("expected: `%s` got: `%s`", tC.want, want)
			}
		})
	}
}
