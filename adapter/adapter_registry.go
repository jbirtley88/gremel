package adapter

import (
	"fmt"
	"io"
	"strings"

	"github.com/jbirtley88/gremel/data"
)

// TODO(jb): not sure about this signature, it might be better to delegate any headings extraction
// to the consumer of the Load() method, so that the adapter can concentrate on loading the data.
// It conceptually doesn't make sense for JSON to have 'headings' but it is necessary when we get
// into the interactive SQL-esque parts of Gremel, so I'm going to leave it here for now.
type AdapterConstructor func(name string, ctx data.GremelContext, source io.ReadCloser) Adapter

type AdapterRegistry interface {
	Register(name string, constructor AdapterConstructor) error
	Get(name string, ctx data.GremelContext, source io.ReadCloser) Adapter
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

func (r *adapterRegistry) Get(name string, ctx data.GremelContext, source io.ReadCloser) Adapter {
	if constructor, found := r.adaptersByName[strings.ToLower(name)]; found {
		return constructor(name, ctx, source)
	}

	return NewErrorAdapter(fmt.Errorf("Adapter '%s' not found", name))
}
