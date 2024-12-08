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
}

type grammar struct {
	tags      map[string]TagCompiler
	blockDefs map[string]*blockSyntax
}

// NewConfig creates a new Settings.
func NewConfig() Config {
	g := grammar{
		tags:      map[string]TagCompiler{},
		blockDefs: map[string]*blockSyntax{},
	}
	return Config{Config: parser.NewConfig(g), grammar: g, Cache: map[string][]byte{}}
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
