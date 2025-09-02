package adapter

import (
	"context"
	"testing"

	"github.com/jbirtley88/gremel/data"
)

func TestCreateDB(t *testing.T) {
	ctx := data.NewGremelContext(context.Background())
	err := CreateDBFromFile(ctx, "gremel", "accounts", "../test_resources/accounts_nested.json")
	if err != nil {
		t.Fatalf("failed to create DB: %v", err)
	}
}
