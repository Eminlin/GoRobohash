package gorobohash

import "testing"

func Test_resource_assemble(t *testing.T) {
	type args struct {
		roboset string
		colors  string
		bgset   string
		format  string
		x       int
		y       int
	}
	tests := []struct {
		name string
		r    *resource
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{},
			r:    NewResource("emin"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.assemble(tt.args.roboset, tt.args.colors, tt.args.bgset, tt.args.format, tt.args.x, tt.args.y)
		})
	}
}
