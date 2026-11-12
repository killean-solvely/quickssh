package ssh

import (
	"testing"
)

func Test_isPubKey(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Is A Pub Key",
			args: args{
				name: "pubkeytest.pub",
			},
			want: true,
		},
		{
			name: "Is Not A Pub Key",
			args: args{
				name: "pubkeytest",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPubKey(tt.args.name); got != tt.want {
				t.Errorf("isPubKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
