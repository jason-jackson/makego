package src

import (
	"strings"
	"testing"
)

func Test_findRouter(t *testing.T) {
	testCases := []struct {
		name    string
		search  string
		want    Router
		wantErr string
	}{
		{
			name:    "empty",
			search:  "",
			wantErr: "no router matching",
		},
		{
			name:   "exact name",
			search: "gin",
			want:   routers["gin"],
		},
		{
			name:   "different case",
			search: "ECHO",
			want:   routers["echo"],
		},
		{
			name:    "not found",
			search:  "does not exist",
			wantErr: "no router matching",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			want, err := findRouter(tC.search)
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
