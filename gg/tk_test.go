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
