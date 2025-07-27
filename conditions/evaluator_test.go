package conditions

import (
	"context"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestVariablesInContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "context_var", 99)
	filter := `{context_var} == 99 AND {x} == 1`
	p := NewParser(strings.NewReader(filter))
	expr, err := p.Parse()
	require.Nil(t, err)

	// Evaluate expression passing data for {vars}
	r, err := Evaluate(ctx, expr, map[string]any{"x": 1})
	require.Nil(t, err)
	require.True(t, r)
}

func TestDeriveZeroValue(t *testing.T) {
	filter := `{x} == 0 AND {y} == 0`
	p := NewParser(strings.NewReader(filter))
	expr, err := p.Parse()
	require.Nil(t, err)

	// Evaluate expression passing data for {vars}
	_ = viper.ReadInConfig()
	viper.Set("config.undefined_is_zero", true)
	require.Nil(t, err)
	r, err := Evaluate(context.TODO(), expr, map[string]any{})
	require.Nil(t, err)
	require.True(t, r)
}
