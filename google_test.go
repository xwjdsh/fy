package fy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoogleTranslate(t *testing.T) {
	{
		resp := GoogleTranslate(context.Background(), Request{
			ToLang: Chinese,
			Text:   "test",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(google)
		expectedResp.Result = "测试"
		assert.Equal(t, expectedResp, resp)
	}

	{
		resp := GoogleTranslate(context.Background(), Request{
			ToLang: English,
			Text:   "测试",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(google)
		expectedResp.Result = "test"
		assert.EqualValues(t, expectedResp, resp)

	}
}
