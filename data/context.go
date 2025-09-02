package data

import (
	"context"
)

type GremelContext interface {
	Context() context.Context
	Values() Metadata
}

type GremelContextImpl struct {
	ctx    context.Context
	values Metadata
}

func NewGremelContext(ctx context.Context, baseline ...Metadata) GremelContext {
	return &GremelContextImpl{
		ctx:    ctx,
		values: NewMetadata(baseline...),
	}
}

func (c *GremelContextImpl) Context() context.Context {
	return c.ctx
}

func (c *GremelContextImpl) Values() Metadata {
	return c.values
}
