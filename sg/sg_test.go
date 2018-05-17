package sg

import (
	"reflect"
	"testing"

	"github.com/xwjdsh/fy"
)

func Test_sogou_Translate(t *testing.T) {
	type args struct {
		req fy.Request
	}
	tests := []struct {
		name     string
		s        *sogou
		args     args
		wantResp *fy.Response
	}{
		{
			name: "text = test",
			s:    &sogou{},
			args: args{
				req: fy.Request{
					TargetLang: "zh-CN",
					Text:       "test",
				},
			},
			wantResp: &fy.Response{
				FullName: "sogou",
				Result:   "试验",
				Err:      nil,
			},
		},
		{
			name: "text = 测试",
			s:    &sogou{},
			args: args{
				req: fy.Request{
					TargetLang: "en",
					Text:       "测试",
				},
			},
			wantResp: &fy.Response{
				FullName: "sogou",
				Result:   "test",
				Err:      nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sogou{}
			if gotResp := s.Translate(tt.args.req); !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("sogou.Translate() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
