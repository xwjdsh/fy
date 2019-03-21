package fy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsChinese(t *testing.T) {
	assert.True(t, IsChinese("测试"))
	assert.False(t, IsChinese("test"))
}
