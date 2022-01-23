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
		{"star", "*", []EntityTag{{true, false, ""}}, false},
		{"starcommas", ",*", []EntityTag{{true, false, ""}}, false},
		{"one", `"wibble"`, []EntityTag{{false, false, "wibble"}}, false},
		{"two", `"wibble","spong"`, []EntityTag{{false, false, "wibble"}, {false, false, "spong"}}, false},
		{"twocommaspace", ` , "wibble" ,, "spong" , `, []EntityTag{{false, false, "wibble"}, {false, false, "spong"}}, false},
		{"weak", `W/"wibble"`, []EntityTag{{false, true, "wibble"}}, false},
		{"weakstrong", `W/"wibble", "spong"`, []EntityTag{{false, true, "wibble"}, {false, false, "spong"}}, false},
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
