package helper

import "testing"

func TestBytesAddUnit(t *testing.T) {
	type args struct {
		bytes int64
	}
	tests := []struct {
		name             string
		args             args
		wantSizeWithUnit string
	}{
		{name: "kb", args: args{bytes: 1200}, wantSizeWithUnit: "1.20 KB"},
		{name: "mb", args: args{bytes: 1200000}, wantSizeWithUnit: "1.20 MB"},
		{name: "gb", args: args{bytes: 1200000000}, wantSizeWithUnit: "1.20 GB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSizeWithUnit := BytesAddUnit(tt.args.bytes); gotSizeWithUnit != tt.wantSizeWithUnit {
				t.Errorf("BytesAddUnit() = %v, want %v", gotSizeWithUnit, tt.wantSizeWithUnit)
			}
		})
	}
}
