package entity

import (
	"testing"

	"github.com/duswie/dxf/format"
)

func TestThreeDSolid_Format(t *testing.T) {

	gen255 := func() string {
		var s string
		for i := 0; i < 255; i++ {
			s += "A"
		}
		for i := 0; i < 255; i++ {
			s += "B"
		}
		return s
	}

	type fields struct {
		entity     *entity
		data_lines []string
		bbox       [][]float64
	}
	type args struct {
		fm format.Formatter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test < 255 characters",
			fields: fields{
				entity:     NewEntity(THREEDSOLID),
				data_lines: []string{"test entry smaller than 255 characters"},
			},
			args: args{
				fm: format.NewASCII(),
			},
		},
		{
			name: "Test > 255 characters",
			fields: fields{
				entity:     NewEntity(THREEDSOLID),
				data_lines: []string{gen255(), gen255()},
			},
			args: args{
				fm: format.NewASCII(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &ThreeDSolid{
				entity:     tt.fields.entity,
				data_lines: tt.fields.data_lines,
				Bbox:       tt.fields.bbox,
			}
			f.Format(tt.args.fm)
			println(tt.args.fm.Output())
		})
	}
}
