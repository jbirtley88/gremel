package data

import (
	"context"
	"fmt"
	"strings"
)

type AdapterConstructor func(name string, ctx context.Context) Adapter

type AdapterRegistry interface {
	Register(name string, constructor AdapterConstructor) error
	Get(name string, ctx context.Context) Adapter
}

var adapterInstance AdapterRegistry = &adapterRegistry{
	adaptersByName: make(map[string]AdapterConstructor),
}

func GetAdapterRegistry() AdapterRegistry {
	return adapterInstance
}

type adapterRegistry struct {
	adaptersByName map[string]AdapterConstructor
}

func (r *adapterRegistry) Register(name string, constructor AdapterConstructor) error {
	r.adaptersByName[strings.ToLower(name)] = constructor
	return nil
}

func (r *adapterRegistry) Get(name string, ctx context.Context) Adapter {
	if constructor, found := r.adaptersByName[strings.ToLower(name)]; found {
		return constructor(name, ctx)
	}

	return NewErrorAdapter(fmt.Errorf("Adapter '%s' not found", name))
}
