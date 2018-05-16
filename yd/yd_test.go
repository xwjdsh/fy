package sg

import (
	"reflect"
	"testing"

	"github.com/xwjdsh/fy"
)

func Test_youdao_Translate(t *testing.T) {
	type args struct {
		req fy.Request
	}
	tests := []struct {
		name     string
		y        *youdao
		args     args
		wantResp *fy.Response
	}{
		{
			name: "text = test",
			y:    &youdao{},
			args: args{
				req: fy.Request{
					TargetLang: "zh-CN",
					Text:       "test",
				},
			},
			wantResp: &fy.Response{
				FullName: "youdao",
				Result:   "测试",
				Err:      nil,
			},
		},
		{
			name: "text = 测试",
			y:    &youdao{},
			args: args{
				req: fy.Request{
					TargetLang: "en",
					Text:       "测试",
				},
			},
			wantResp: &fy.Response{
				FullName: "youdao",
				Result:   "test",
				Err:      nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &youdao{}
			if gotResp := y.Translate(tt.args.req); !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("youdao.Translate() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
