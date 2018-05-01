package gg

import "testing"

func Test_calTK(t *testing.T) {
	type args struct {
		vq    string
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{
				vq:    "423665.134779550",
				query: "test",
			},
			want: "238118.382167",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calTK(tt.args.vq, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("calTK() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calTK() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getVq(t *testing.T) {
	type args struct {
		dataStr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{
				dataStr: `dataStr = ` + `LOW_CONFIDENCE_THRESHOLD=-1;MAX_ALTERNATIVES_ROUNDTRIP_RESULTS=1;TKK=eval('((function(){var a\x3d1966732470;var b\x3d1714107181;return 423123+\x27.\x27+(a+b)})())');VERSION_LABEL = 'twsfe_w_20180402_RC00';`,
			},
			want: "423123.3680839651",
		},
		{
			args:    args{dataStr: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getVq(tt.args.dataStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getVq() = %v, want %v", got, tt.want)
			}
		})
	}
}
