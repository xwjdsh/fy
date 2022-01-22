package fy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeeplTranslate(t *testing.T) {
	{
		resp := DeeplTranslate(context.Background(), Request{
			ToLang: Chinese,
			Text:   "test",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(deepl)
		expectedResp.Result = "测试"
		assert.Equal(t, expectedResp, resp)
	}

	{
		resp := DeeplTranslate(context.Background(), Request{
			ToLang: English,
			Text:   "测试",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(deepl)
		expectedResp.Result = "test (machinery etc)"
		assert.Equal(t, expectedResp, resp)
	}
}
