package markdown

import (
	"reflect"
	"testing"
)

func Test_getAllMatchMD(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantPaths []string
	}{
		// TODO: Add test cases.
		{name: "1", path: "./*.md", wantPaths: []string{"ass"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPaths := getAllMatchMD(tt.path)
			if !reflect.DeepEqual(gotPaths, tt.wantPaths) {
				t.Errorf("getAllMatchMD() = %v, want %v", gotPaths, tt.wantPaths)
			}
		})
	}
}
