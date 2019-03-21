package fy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSogouTranslate(t *testing.T) {
	{
		resp := SogouTranslate(context.Background(), Request{
			ToLang: Chinese,
			Text:   "test",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(sogou)
		expectedResp.Result = "试验"
		assert.Equal(t, expectedResp, resp)
	}

	{
		resp := SogouTranslate(context.Background(), Request{
			ToLang: English,
			Text:   "测试",
		})
		assert.Nil(t, resp.Err)
		expectedResp := newResp(sogou)
		expectedResp.Result = "test"
		assert.Equal(t, expectedResp, resp)
	}
}
