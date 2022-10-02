package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	str := "Apple"
	want := "elppA"
	got := Reverse(str)
	assert.Equalf(t, want, got, "Got %s but want %s", got, want)
}

func TestLookupEnvOrDefault(t *testing.T) {
	//set some values
	os.Setenv("BAR", "Not Bar")
	testCases := map[string]struct {
		name       string
		defaultVal string
		want       string
	}{
		"nilDefaults": {
			name:       "FOO",
			defaultVal: "",
			want:       "",
		},
		"defaults": {
			name:       "FOO",
			defaultVal: "foo",
			want:       "foo",
		},
		"nodefaults": {
			name:       "BAR",
			defaultVal: "bar",
			want:       "Not Bar",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := LookupEnvOrString(tc.name, tc.defaultVal)
			assert.Equalf(t, tc.want, got, "Got %s but want %s", got, tc.want)
		})
	}
}
