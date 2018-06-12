package util

import (
	"reflect"
	"testing"
)

func TestParseEntityTags(t *testing.T) {
	tests := []struct {
		name     string
		args     string
		wantTags []EntityTag
		wantErr  bool
	}{
		{"empty", "", nil, true},
		{"spaces", "  ", nil, true},
		{"commas", ",,", nil, true},
		{"star", "*", []EntityTag{EntityTag{true, false, ""}}, false},
		{"starcommas", ",*", []EntityTag{EntityTag{true, false, ""}}, false},
		{"one", `"wibble"`, []EntityTag{EntityTag{false, false, "wibble"}}, false},
		{"two", `"wibble","spong"`, []EntityTag{EntityTag{false, false, "wibble"}, EntityTag{false, false, "spong"}}, false},
		{"twocommaspace", ` , "wibble" ,, "spong" , `, []EntityTag{EntityTag{false, false, "wibble"}, EntityTag{false, false, "spong"}}, false},
		{"weak", `W/"wibble"`, []EntityTag{EntityTag{false, true, "wibble"}}, false},
		{"weakstrong", `W/"wibble", "spong"`, []EntityTag{EntityTag{false, true, "wibble"}, EntityTag{false, false, "spong"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTags, err := ParseEntityTags(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEntityTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTags, tt.wantTags) {
				t.Errorf("ParseEntityTags() = %v, want %v", gotTags, tt.wantTags)
			}
		})
	}
}
