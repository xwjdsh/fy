package bd

import (
	"reflect"
	"testing"

	"github.com/xwjdsh/fy"
)

func Test_baidu_Translate(t *testing.T) {
	type args struct {
		req *fy.Request
	}
	tests := []struct {
		name     string
		b        *baidu
		args     args
		wantResp *fy.Response
	}{
		{
			name: "text = test",
			b:    &baidu{},
			args: args{
				req: &fy.Request{
					IsChinese: false,
					Text:      "test",
				},
			},
			wantResp: &fy.Response{
				FullName: "baidu",
				Result:   "测试",
				Err:      nil,
			},
		},
		{
			name: "text = 测试",
			b:    &baidu{},
			args: args{
				req: &fy.Request{
					IsChinese: true,
					Text:      "测试",
				},
			},
			wantResp: &fy.Response{
				FullName: "baidu",
				Result:   "test",
				Err:      nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &baidu{}
			if gotResp := b.Translate(tt.args.req); !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("baidu.Translate() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
