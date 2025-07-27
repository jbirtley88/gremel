package data

import (
	"context"
)

// An Adapter must be cabale of doing two things:
//
// 1. Fetch data from a remote source, such as an HTTP endpoint.
// 2. Provide a string representation of the data that can be displayed through the pseudo-filesystem.
type Adapter interface {
	GetName() string
	Fetch() error
}

type BaseAdapter struct {
	Name string
	Ctx  context.Context
}

func NewBaseAdapter(name string, ctx context.Context) *BaseAdapter {
	a := &BaseAdapter{
		Name: name,
		Ctx:  ctx,
	}
	if a.Ctx == nil {
		a.Ctx = context.TODO()
	}
	return a
}

func (a *BaseAdapter) GetName() string {
	return a.Name
}

func (a *BaseAdapter) Filter(criteria string) error {
	return nil
}
