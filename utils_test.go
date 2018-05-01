package fy

import "testing"

func TestIsChinese(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "str = test",
			args: args{
				str: "test",
			},
			want: false,
		},
		{
			name: "str = 测试",
			args: args{
				str: "测试",
			},
			want: true,
		},
		{
			name: "str = test 测试",
			args: args{
				str: "test 测试",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsChinese(tt.args.str); got != tt.want {
				t.Errorf("IsChinese() = %v, want %v", got, tt.want)
			}
		})
	}
}
