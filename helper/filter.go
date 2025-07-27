package helper

import (
	"context"
	"log"
	"strings"

	"github.com/jbirtley88/gremel/conditions"
)

func Matches(ctx context.Context, row map[string]any, filter string) (bool, error) {
	if filter == "" {
		return true, nil
	}

	// Parse the condition language and get expression
	p := conditions.NewParser(strings.NewReader(filter))
	expr, err := p.Parse()
	if err != nil {
		log.Fatal("Matches(): '%s':", filter, err)
	}

	// Evaluate expression passing data for {vars}
	r, err := conditions.Evaluate(ctx, expr, row)
	if err != nil {
		log.Fatal("Matches(): '%s':", filter, err)
	}

	return r, nil
}
