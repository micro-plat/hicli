package tmpl

import (
	"os"
	"reflect"
	"testing"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		append   bool
		wantFile *os.File
		wantErr  bool
	}{
		// TODO: Add test cases.
		{name: "1", path: "./conf.go", append: true, wantFile: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFile, err := Create(tt.path, tt.append)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFile, tt.wantFile) {
				t.Errorf("Create() = %v, want %v", gotFile, tt.wantFile)
			}
		})
	}
}
