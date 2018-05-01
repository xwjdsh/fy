package qq

import (
	"reflect"
	"testing"

	"github.com/xwjdsh/fy"
)

func Test_tencent_Translate(t *testing.T) {
	type args struct {
		req *fy.Request
	}
	tests := []struct {
		name     string
		t        *tencent
		args     args
		wantResp *fy.Response
	}{
		{
			name: "text = test",
			t:    &tencent{},
			args: args{
				req: &fy.Request{
					IsChinese: false,
					Text:      "test",
				},
			},
			wantResp: &fy.Response{
				FullName: "tencent",
				Result:   "试验 / 测验 / 化验 / 检查",
				Err:      nil,
			},
		},
		{
			name: "text = 测试",
			t:    &tencent{},
			args: args{
				req: &fy.Request{
					IsChinese: true,
					Text:      "测试",
				},
			},
			wantResp: &fy.Response{
				FullName: "tencent",
				Result:   "test",
				Err:      nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			te := &tencent{}
			if gotResp := te.Translate(tt.args.req); !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("tencent.Translate() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_getQtk(t *testing.T) {
	type args struct {
		dataStr string
	}
	tests := []struct {
		name    string
		args    args
		wantQtv string
		wantQtk string
		wantErr bool
	}{
		{
			args: args{
				dataStr: `
					document.cookie = "qtv=ad15088b8bcde724";
				  document.cookie = "qtk=aK4qrfL4bLogktVEfIMc785lhWKxOuLuOF243HgKs/lOcPqPhoiwsR+7ysGoTF/rqx1EABKUpNJq2OqbE1PY9T9ICiU2Qm2l0yIMqg3mworcjCX8tiaZzYjkQQqSTk7gCIz/WY0NhTJUrrOemb6nRQ==";
				`,
			},
			wantQtv: "ad15088b8bcde724",
			wantQtk: "aK4qrfL4bLogktVEfIMc785lhWKxOuLuOF243HgKs/lOcPqPhoiwsR+7ysGoTF/rqx1EABKUpNJq2OqbE1PY9T9ICiU2Qm2l0yIMqg3mworcjCX8tiaZzYjkQQqSTk7gCIz/WY0NhTJUrrOemb6nRQ==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQtv, gotQtk, err := getQtk(tt.args.dataStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("getQtk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotQtv != tt.wantQtv {
				t.Errorf("getQtk() gotQtv = %v, want %v", gotQtv, tt.wantQtv)
			}
			if gotQtk != tt.wantQtk {
				t.Errorf("getQtk() gotQtk = %v, want %v", gotQtk, tt.wantQtk)
			}
		})
	}
}
