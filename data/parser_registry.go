package data

import (
	"context"
	"fmt"
	"strings"
)

type ParserConstructor func(name string, ctx context.Context) Parser

type ParserRegistry interface {
	Register(name string, constructor ParserConstructor) error
	Get(name string, ctx context.Context) Parser
}

var parserInstance ParserRegistry = &parserRegistry{
	parsersByName: make(map[string]ParserConstructor),
}

func GetParserRegistry() ParserRegistry {
	return parserInstance
}

type parserRegistry struct {
	parsersByName map[string]ParserConstructor
}

func (r *parserRegistry) Register(name string, constructor ParserConstructor) error {
	r.parsersByName[strings.ToLower(name)] = constructor
	return nil
}

func (r *parserRegistry) Get(name string, ctx context.Context) Parser {
	if constructor, found := r.parsersByName[strings.ToLower(name)]; found {
		return constructor(name, ctx)
	}

	return NewParserError(fmt.Errorf("Parser '%s' not found", name))
}
