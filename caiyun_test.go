package fy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaiyunTranslate(t *testing.T) {
	{
		resp := CaiyunTranslate(context.Background(), Request{
			ToLang: Chinese,
			Text:   "One",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(caiyun)
		expectedResp.Result = "一个"
		assert.Equal(t, expectedResp, resp)
	}

	{
		resp := CaiyunTranslate(context.Background(), Request{
			ToLang: English,
			Text:   "一",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(caiyun)
		expectedResp.Result = "One"
		assert.Equal(t, expectedResp, resp)
	}
}
