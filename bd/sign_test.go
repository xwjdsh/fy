package bd

import "testing"

func Test_calSign(t *testing.T) {
	type args struct {
		gtk   string
		query string
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
				gtk:   "320305.131321201",
				query: "test",
			},
			want: "431039.159886",
		},
		{
			name: "测试",
			args: args{
				gtk:   "320305.131321201",
				query: "测试",
			},
			want: "706553.926920",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calSign(tt.args.gtk, tt.args.query)
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

func Test_getTokenAndGtk(t *testing.T) {
	type args struct {
		dataStr string
	}
	tests := []struct {
		name      string
		args      args
		wantToken string
		wantGtk   string
		wantErr   bool
	}{
		{
			args:    args{dataStr: ""},
			wantErr: true,
		},
		{
			args: args{
				dataStr: `
				<script>window.bdstoken = '';window.gtk = '320305.131321201';</script>
<script type="text/javascript" src="//fanyi.bdstatic.com/static/translation/lib/third_party/mod_b80b8f9.js"></script>
<script type="text/javascript" src="//fanyi.bdstatic.com/static/translation/pkg/public_c587d32.js"></script>
<script type="text/javascript" src="//fanyi.bdstatic.com/static/translation/pkg/index_01eb09e.js"></script>

				window['common'] = {
    token: '395c9e69e667bf7579a2a3ade6391edb',
    systime: '1525177203477',
    logid: 'a831f7926cd2a9bd8078cf7968894a01',
				`,
			},
			wantToken: "395c9e69e667bf7579a2a3ade6391edb",
			wantGtk:   "320305.131321201",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, gotGtk, err := getTokenAndGtk(tt.args.dataStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTokenAndGtk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("getTokenAndGtk() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
			if gotGtk != tt.wantGtk {
				t.Errorf("getTokenAndGtk() gotGtk = %v, want %v", gotGtk, tt.wantGtk)
			}
		})
	}
}
