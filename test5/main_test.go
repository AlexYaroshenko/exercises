package main

import "testing"

func Test_nonRepeating(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Case 1",
			args: args{s: "swiss"},
			want: "w",
		},
		{
			name: "Test Case 2",
			args: args{s: "racecar"},
			want: "e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nonRepeating(tt.args.s); got != tt.want {
				t.Errorf("nonRepeating() = %v, want %v", got, tt.want)
			}
		})
	}
}
