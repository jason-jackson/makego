package src

import (
	"strings"
	"testing"
)

func Test_findDatabase(t *testing.T) {
	testCases := []struct {
		name    string
		search  string
		want    Database
		wantErr string
	}{
		{
			name:    "empty",
			search:  "",
			wantErr: "no database matching",
		},
		{
			name:   "exact name",
			search: "mysql",
			want:   databases["mysql"],
		},
		{
			name:   "exact match",
			search: "pg",
			want:   databases["postgres"],
		},
		{
			name:   "different case",
			search: "postgreSQL",
			want:   databases["postgres"],
		},
		{
			name:    "not found",
			search:  "does not exist",
			wantErr: "no database matching",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			want, err := findDatabase(tC.search)
			if tC.wantErr != "" {
				if !strings.Contains(err.Error(), tC.wantErr) {
					t.Errorf("expected `%s` to contain `%s`", err.Error(), tC.wantErr)
				}
				return
			}
			if tC.want.Name != want.Name {
				t.Errorf("expected: `%s` got: `%s`", tC.want.Name, want.Name)
			}
		})
	}
}
