package render

import (
	"github.com/chrisatbd/liquid/parser"
)

// Config holds configuration information for parsing and rendering.
type Config struct {
	parser.Config
	grammar
	Cache           map[string][]byte
	StrictVariables bool
	TemplateStore   ITemplateStore
}

type grammar struct {
	tags      map[string]TagCompiler
	blockDefs map[string]*blockSyntax
}

// NewConfig creates a new Settings.  TemplateStore is initialized to a FileTemplateStore for backwards compatability
func NewConfig() Config {
	g := grammar{
		tags:      map[string]TagCompiler{},
		blockDefs: map[string]*blockSyntax{},
	}
	return Config{
		Config:        parser.NewConfig(g),
		grammar:       g,
		Cache:         map[string][]byte{},
		TemplateStore: &FileTemplateStore{},
	}
}
