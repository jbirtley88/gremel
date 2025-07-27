package data

import (
	"fmt"
	"io"
)

// ErrorAdapter is what gets returned by the adpater registry if there is any
// problem encountered when trying to instantiate an adapter.
//
// It is fully compliant with the Adapter interface, so it can be returned by
// adapter constructors, but all of its methods report / return an error.
//
// The most common use-case for this is an attempt to look up an adapter
// which has not been configured in config.yml, e.g.
//
//	curl http://.../proxy/the_adapter
//	                      ^^^^^^^^^^^
//	                      there is no authfairy.adapters.the_adapter in config.yml
type ErrorAdapter struct {
	BaseAdapter
	Err error
}

func NewErrorAdapter(err error) Adapter {
	return &ErrorAdapter{
		Err: err,
	}
}

func (a *ErrorAdapter) GetName() string {
	return fmt.Sprintf("ERROR: %v", a.Err)
}

func (a *ErrorAdapter) Fetch() error {
	return fmt.Errorf("Fetch(): %v", a.Err)
}

func (a *ErrorAdapter) ResponseString(content []byte, writeTo io.Writer) error {
	_, err := writeTo.Write([]byte(a.Err.Error()))
	return err
}
