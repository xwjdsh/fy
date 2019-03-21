package fy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTencentTranslate(t *testing.T) {
	{
		resp := TencentTranslate(context.Background(), Request{
			ToLang: Chinese,
			Text:   "test",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(tencent)
		expectedResp.Result = "试验 / 测试 / 烤钵 / 检验"
		assert.Equal(t, expectedResp, resp)
	}

	{
		resp := TencentTranslate(context.Background(), Request{
			ToLang: English,
			Text:   "测试",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(tencent)
		expectedResp.Result = "test"
		assert.Equal(t, expectedResp, resp)
	}

}
