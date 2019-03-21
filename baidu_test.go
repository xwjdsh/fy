package fy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaiduTranslate(t *testing.T) {
	{
		resp := BaiduTranslate(context.Background(), Request{
			ToLang: Chinese,
			Text:   "test",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(baidu)
		expectedResp.Result = "测试"
		assert.Equal(t, expectedResp, resp)
	}

	{
		resp := BaiduTranslate(context.Background(), Request{
			ToLang: English,
			Text:   "测试",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(baidu)
		expectedResp.Result = "test"
		assert.Equal(t, expectedResp, resp)
	}
}
