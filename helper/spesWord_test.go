package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeplaceAllSpecialChars(t *testing.T) { // MTY0MDA5MDUwMQ
	got := PeplaceAllSpecialChars(`{
				"id":"@UUID",
				"time":"@Unix",
				"num":"@Int",
				"content":"@Base64",
				"phone":"@Phone",
				"describe":"test replace specail chars",
				"end":"end-@UUID",
			}`)
	t.Logf("got %+v", got)
}

func TestCalcNewStringLength(t *testing.T) {
	type args struct {
		old string
	}
	tests := []struct {
		name             string
		args             args
		wantNewLen       int
		wantReplaceCount int
	}{
		{name: "empty", args: args{old: ""}, wantNewLen: 0, wantReplaceCount: 0},
		{
			name:             "base",
			args:             args{old: `{"id":"@UUID","time":"@Unix","num":"@Int","phone":"@Phone","describe":"test replace specail chars","end":"end-@UUID"}`},
			wantNewLen:       117 + 31 + 5 + 6 + 31 + 5,
			wantReplaceCount: 5,
		},
		{
			name:             "Boundary-end",
			args:             args{old: `{"end":"end-@UUI"}`},
			wantNewLen:       18,
			wantReplaceCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewLen, gotReplaceCount := CalcNewStringLength(tt.args.old)
			assert.Equal(t, tt.wantNewLen, gotNewLen)
			assert.Equal(t, tt.wantReplaceCount, gotReplaceCount)
		})
	}
}

func TestIndexSpecialChars(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name      string
		args      args
		wantIndex int
		wantM     Mocker
	}{
		{name: "empty", args: args{data: ""}, wantIndex: -1, wantM: nil},
		{name: "no Match", args: args{data: "data:123"}, wantIndex: -1, wantM: nil},
		{name: "interage", args: args{data: "data:@Int"}, wantIndex: 5, wantM: new(Int)},
		{name: "unix", args: args{data: "data:@Unix"}, wantIndex: 5, wantM: new(Unix)},
		{name: "uuid", args: args{data: "data:@UUID"}, wantIndex: 5, wantM: new(UUID)},
		{name: "uuid", args: args{data: "data:@Phone"}, wantIndex: 5, wantM: new(Phone)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIndex, gotM := IndexSpecialChars(tt.args.data)
			assert.EqualValues(t, tt.wantM, gotM)
			assert.Equal(t, tt.wantIndex, gotIndex)
		})
	}
}
