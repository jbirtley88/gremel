package helper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCommandOutputNoTimeout(t *testing.T) {
	// Sleep 1 for 1 second
	cmdAndArgs := []string{"sh", "-c", `echo "this is stdout" ; sleep 1 ; echo "this is stderr" 1>&2`}
	// Timeout is 2 seconds
	result, err := GetCommandOutput(cmdAndArgs, time.Second*2)
	require.Nil(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 0, result.ExitCode)
	assert.False(t, result.IsTimeout)
	assert.Equal(t, "this is stdout", string(result.Stdout))
	assert.Equal(t, `this is stderr`, string(result.Stderr))
}

func TestGetCommandOutputWithTimeout(t *testing.T) {
	// Sleep 1 for 1 second
	cmdAndArgs := []string{"sh", "-c", `echo "this is stdout" ; sleep 2 ; echo "this is stderr" 1>&2`}
	// Timeout is 2 seconds
	result, err := GetCommandOutput(cmdAndArgs, time.Second*1)
	require.Nil(t, err)
	require.NotNil(t, result)
	assert.Equal(t, -1, result.ExitCode)
	assert.True(t, result.IsTimeout)
	assert.Equal(t, "this is stdout", string(result.Stdout))
	// The 'echo this is stderr' command will not have executed due to the timeout
	assert.Equal(t, ``, string(result.Stderr))
}
