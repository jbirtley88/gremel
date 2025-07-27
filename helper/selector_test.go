package helper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseSelector(t *testing.T) {
	selector := NewBaseSelector(context.TODO())
	assert.NotNil(t, selector)
	assert.Equal(t, "default", selector.GetName())
}

func TestBaseSelectorFilter(t *testing.T) {
	selector := NewBaseSelector(context.TODO())
	input := []map[string]any{
		{"key1": "value1"},
		{"key2": "value2"},
	}
	result, err := selector.Where(input, "")
	require.Nil(t, err)
	assert.Equal(t, input, result)
}

func TestBaseSelectorOrder(t *testing.T) {
	selector := NewBaseSelector(context.TODO())
	input := []map[string]any{
		{"key1": "value1"},
		{"key2": "value2"},
	}
	result, err := selector.Order(input, "")
	require.Nil(t, err)
	assert.Equal(t, input, result)
}
