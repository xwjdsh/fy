package fy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYoudaoTranslate(t *testing.T) {
	{
		resp := YoudaoTranslate(context.Background(), Request{
			ToLang: Chinese,
			Text:   "test",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(youdao)
		expectedResp.Result = "测试"
		assert.Equal(t, expectedResp, resp)
	}

	{
		resp := YoudaoTranslate(context.Background(), Request{
			ToLang: English,
			Text:   "测试",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(youdao)
		expectedResp.Result = "test"
		assert.Equal(t, expectedResp, resp)
	}
}
