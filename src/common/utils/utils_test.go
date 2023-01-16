package utils

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var currentDir, _ = os.Getwd()

func TestGetFileFormat(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "yaml",
			args: args{filePath: "dir/file.yaml"},
			want: "yaml",
		},
		{
			name: "yml",
			args: args{filePath: "dir/file.yml"},
			want: "yml",
		},
		{
			name: "json",
			args: args{filePath: "dir/file.json"},
			want: "json",
		},
		{
			name: "no file type",
			args: args{filePath: "dir/file"},
			want: "",
		},
		{
			name: "empty string",
			args: args{filePath: ""},
			want: "",
		},
		{
			name: "template-yaml",
			args: args{filePath: currentDir + "/../../../tests/cloudformation/resources/extensions/ebs.template"},
			want: "yaml",
		},
		{
			name: "template-yaml",
			args: args{filePath: currentDir + "/../../../tests/cloudformation/resources/extensions/ebs2.template"},
			want: "json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileFormat(tt.args.filePath); got != tt.want {
				t.Errorf("GetFileFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInSlice(t *testing.T) {
	type args struct {
		slice interface{}
		elem  interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "in slice string",
			args: args{slice: []string{"a", "b", "c", "e"}, elem: "a"},
			want: true,
		},
		{
			name: "not in slice string",
			args: args{slice: []string{"a", "b", "c", "e"}, elem: "d"},
			want: false,
		},
		{
			name: "in slice int",
			args: args{slice: []int{1, 2, 3, 4}, elem: 1},
			want: true,
		},
		{
			name: "not in slice int",
			args: args{slice: []int{1, 2, 3, 4}, elem: 5},
			want: false,
		},
		{
			name: "slice in slice ",
			args: args{slice: [][]int{{1, 2, 3, 4}, {5, 6}, {7}}, elem: []int{5, 6}},
			want: true,
		},
		{
			name: "not slice in slice ",
			args: args{slice: [][]int{{1, 2, 3, 4}, {5, 6}, {7}}, elem: []int{5, 7}},
			want: false,
		},
		{
			name: "different kinds",
			args: args{slice: []int{1, 2, 3, 4}, elem: "bana"},
			want: false,
		},
		{
			name: "nil slice",
			args: args{slice: nil, elem: "bana"},
			want: false,
		},
		{
			name: "empty slice",
			args: args{slice: []int{}, elem: "bana"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InSlice(tt.args.slice, tt.args.elem); got != tt.want {
				t.Errorf("InSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinInt(t *testing.T) {
	t.Run("Test MinInt", func(t *testing.T) {
		ans := MinInt(3, 4)
		assert.Equal(t, 3, ans)
	})
}

func TestSplitStringByComma(t *testing.T) {
	tests := []struct {
		name string
		arg  []string
		want int
	}{
		{
			name: "Test no comma",
			arg:  []string{"Hello"},
			want: 1,
		},
		{
			name: "Test actual list",
			arg:  []string{"Hello", "World"},
			want: 2,
		},
		{
			name: "Test comma delimited list",
			arg:  []string{"tests,.git,node_modules"},
			want: 3,
		},
		{
			name: "Test combined",
			arg:  []string{"tests,.git,node_modules", ".github"},
			want: 4,
		},
	}
	for _, tt := range tests {
		got := len(SplitStringByComma(tt.arg))
		assert.Equal(t, tt.want, got, fmt.Sprintf("Expected to get %v, got %v for \"%v\"", tt.want, got, tt.name))
	}
}

func TestGetEnv(t *testing.T) {
	t.Run("TestExistingEnvVar", func(t *testing.T) {
		_ = os.Setenv("test", "20")
		assert.Equal(t, "20", GetEnv("test", "1"))
		_ = os.Unsetenv("test")
	})

	t.Run("TestExistingEnvVar", func(t *testing.T) {
		_ = os.Setenv("test2", "20")
		assert.Equal(t, "1", GetEnv("test", "1"))
		_ = os.Unsetenv("test2")
	})
}

func TestAllNil(t *testing.T) {
	t.Run("TestCheckForInterfaceWithString", func(t *testing.T) {
		var i interface{}
		i = []string{"bla"}
		assert.Equal(t, false, AllNil(i))
	})
	t.Run("TestCheckForNonInterfaceWithStringIsNotCrashing", func(t *testing.T) {
		var i interface{}
		i = nil
		assert.Equal(t, true, AllNil(i))
	})
	t.Run("TestCheckForInterfaceWithEmptyString", func(t *testing.T) {
		var i interface{}
		i = []interface{}(nil)
		assert.Equal(t, true, AllNil(i))
	})
}
