package fy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBingTranslate(t *testing.T) {
	{
		resp := BingTranslate(context.Background(), Request{
			ToLang: Chinese,
			Text:   "test",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(bing)
		expectedResp.Result = "测试"
		assert.Equal(t, expectedResp, resp)
	}

	{
		resp := BingTranslate(context.Background(), Request{
			ToLang: English,
			Text:   "测试",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(bing)
		expectedResp.Result = "Test"
		assert.Equal(t, expectedResp, resp)
	}
}
