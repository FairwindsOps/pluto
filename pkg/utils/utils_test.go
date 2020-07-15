package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringInSlice(t *testing.T) {
	type args struct {
		a    string
		list []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test true",
			args: args{
				a:    "string",
				list: []string{"test", "string"},
			},
			want: true,
		},
		{
			name: "test false",
			args: args{
				a:    "string",
				list: []string{"test", "nothere"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInSlice(tt.args.a, tt.args.list); got != tt.want {
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestIsFileOrStdin(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		// TODO: Add test cases.
		{"stdin", "-", true},
		{"no file", "notafile.foo", false},
		{"file", "utils.go", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsFileOrStdin(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
