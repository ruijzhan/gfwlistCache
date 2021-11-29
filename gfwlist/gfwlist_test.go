package gfwlist

import "testing"

func Test_gfwList_has(t *testing.T) {
	gfwlist := new()
	type args struct {
		domain string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				domain: "youtube.com",
			},
			want: true,
		},
		{
			args: args{
				domain: "www.youtube.com",
			},
			want: true,
		},
		{
			args: args{
				domain: "www.google.com.hk",
			},
			want: true,
		},
		{
			args: args{
				domain: "www.163.com",
			},
			want: false,
		},
		{
			args: args{
				domain: "www.google.com",
			},
			want: true,
		},
		{
			args: args{
				domain: "www.ox.ac.uk",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gfwlist
			if got := g.has(tt.args.domain); got != tt.want {
				t.Errorf("gfwList.has() = %v, want %v", got, tt.want)
			}
		})
	}
}
