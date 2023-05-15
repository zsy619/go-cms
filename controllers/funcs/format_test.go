package funcs

import "testing"

func TestSizeFormat(t *testing.T) {
	type args struct {
		size      string
		delimiter string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{"1024", ""}, "1KB"},
		{"2", args{"1024", " "}, "1 KB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SizeFormat(tt.args.size, tt.args.delimiter); got != tt.want {
				t.Errorf("SizeFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
