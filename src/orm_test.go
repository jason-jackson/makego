package src

import (
	"strings"
	"testing"
)

func Test_findORM(t *testing.T) {
	testCases := []struct {
		name    string
		search  string
		want    ORM
		wantErr string
	}{
		{
			name:    "empty",
			search:  "",
			wantErr: "no ORM matching",
		},
		{
			name:   "exact name",
			search: "gorm",
			want:   orms["gorm"],
		},
		// {
		// 	name:   "exact match",
		// 	search: "gorm",
		// 	want:   databases["gorm"],
		// },
		{
			name:   "different case",
			search: "GORM",
			want:   orms["gorm"],
		},
		{
			name:    "not found",
			search:  "does not exist",
			wantErr: "no ORM matching",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			want, err := findORM(tC.search)
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
