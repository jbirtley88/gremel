package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHttpHelperBuilderBuildReturnsDifferentInstances(t *testing.T) {
	builder := NewHttpHelperBuilder()

	instance1 := builder.Build()
	instance2 := builder.Build()

	require.NotNil(t, instance1)
	require.NotNil(t, instance2)

	// Verify that the two instances are different objects in memory
	assert.NotSame(t, instance1, instance2, "Build() should return different instances")

	// Both should implement HttpHelper interface
	var _ HttpHelper = instance1
	var _ HttpHelper = instance2
}
