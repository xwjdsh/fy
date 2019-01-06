package sg

import "testing"

func Test_calSign(t *testing.T) {
	type args struct {
		from string
		to   string
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				from: "auto",
				to:   "zh-CHS",
				text: "test",
			},
			want: "f4919df0708ff951580adfde0085a4cc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calSign(tt.args.from, tt.args.to, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("calSign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calSign() = %v, want %v", got, tt.want)
			}
		})
	}
}
