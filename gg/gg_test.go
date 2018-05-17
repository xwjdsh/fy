package gg

import (
	"reflect"
	"testing"

	"github.com/xwjdsh/fy"
)

func Test_google_Translate(t *testing.T) {
	type args struct {
		req fy.Request
	}
	tests := []struct {
		name     string
		g        *google
		args     args
		wantResp *fy.Response
	}{
		{
			name: "text = test",
			g:    &google{},
			args: args{
				req: fy.Request{
					TargetLang: "zh-CN",
					Text:       "test",
				},
			},
			wantResp: &fy.Response{
				FullName: "google",
				Result:   "测试",
				Err:      nil,
			},
		},
		{
			name: "text = 测试",
			g:    &google{},
			args: args{
				req: fy.Request{
					TargetLang: "en",
					Text:       "测试",
				},
			},
			wantResp: &fy.Response{
				FullName: "google",
				Result:   "test",
				Err:      nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &google{}
			if gotResp := s.Translate(tt.args.req); !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("google.Translate() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
