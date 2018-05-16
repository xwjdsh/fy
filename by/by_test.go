package by

import (
	"reflect"
	"testing"

	"github.com/xwjdsh/fy"
)

func Test_bing_Translate(t *testing.T) {
	type args struct {
		req fy.Request
	}
	tests := []struct {
		name     string
		b        *bing
		args     args
		wantResp *fy.Response
	}{
		{
			name: "text = test",
			b:    new(bing),
			args: args{
				req: fy.Request{
					TargetLang: "zh-CN",
					Text:       "test",
				},
			},
			wantResp: &fy.Response{
				FullName: "bing",
				Result:   "测试",
				Err:      nil,
			},
		},
		{
			name: "text = 测试",
			b:    new(bing),
			args: args{
				req: fy.Request{
					TargetLang: "en",
					Text:       "测试",
				},
			},
			wantResp: &fy.Response{
				FullName: "bing",
				Result:   "Test",
				Err:      nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bing{}
			if gotResp := b.Translate(tt.args.req); !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("bing.Translate() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
