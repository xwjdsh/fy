package fy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaidu(t *testing.T) {
	c := New()
	bd := c.tm[BAIDU]
	{
		resp := c.Baidu(context.Background(), "test")
		assert.Nil(t, resp.Err)
		expectedResp := &Response{
			Name:     bd.name(),
			Homepage: bd.homepage(),
			Result:   "测试",
		}
		assert.Equal(t, expectedResp, resp)
	}

	{
		resp := c.Baidu(context.Background(), "测试")
		assert.Nil(t, resp.Err)
		expectedResp := &Response{
			Name:     bd.name(),
			Homepage: bd.homepage(),
			Result:   "test",
		}
		assert.Equal(t, expectedResp, resp)
	}
}

func TestBing(t *testing.T) {
	c := New()
	b := c.tm[BING]

	{
		resp := c.Bing(context.Background(), "test")
		assert.Nil(t, resp.Err)
		expectedResp := &Response{
			Name:     b.name(),
			Homepage: b.homepage(),
			Result:   "测试",
		}
		assert.Equal(t, expectedResp, resp)
	}

	{
		resp := c.Bing(context.Background(), "测试")
		assert.Nil(t, resp.Err)
		expectedResp := &Response{
			Name:     b.name(),
			Homepage: b.homepage(),
			Result:   "test",
		}
		assert.Equal(t, expectedResp, resp)
	}
}

func TestIsChinese(t *testing.T) {
	assert.True(t, isChinese("测试"))
	assert.False(t, isChinese("test"))
}
