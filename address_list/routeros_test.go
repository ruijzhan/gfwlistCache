package addresslist

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		apiAddr string
		user    string
		passwd  string
	}
	tests := []struct {
		name string
		args args
		want AddressList
	}{
		{
			args: args{
				apiAddr: "192.168.1.1:8728",
				user:    "admin",
				passwd:  "aca04rz.",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.args.apiAddr, tt.args.user, tt.args.passwd)
			for !l.Synced() {
				t.Log("Wait for sync...")
				time.Sleep(time.Second)
			}
			t.Log(l.Has("www.google.com"))
			t.Log(l.Has("www.youtube.com"))
			err := l.Add("www.163.com", "1d")
			if err != nil {
				if err != ErrAlreadyHaveSuchEntry {
					t.Fatal(err)
				}
			}
			l.Stop()
		})
	}
}
