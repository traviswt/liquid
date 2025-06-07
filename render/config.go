package render

import (
	"github.com/osteele/liquid/parser"
	"sort"
)

// Config holds configuration information for parsing and rendering.
type Config struct {
	parser.Config
	grammar
	Cache           map[string][]byte
	StrictVariables bool
	TemplateStore   TemplateStore
}

type grammar struct {
	tags      map[string]TagCompiler
	blockDefs map[string]*blockSyntax
}

// NewConfig creates a new Settings.
// TemplateStore is initialized to a FileTemplateStore for backwards compatibility
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

func (c Config) ListFilters() []string {
	return c.Config.ListFilters()
}

func (c Config) ListTags() []string {
	var l []string
	for k := range c.tags {
		l = append(l, k)
	}
	sort.Strings(l)
	return l
}
