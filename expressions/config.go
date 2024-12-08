package expressions

import "sort"

// Config holds configuration information for expression interpretation.
type Config struct {
	filters map[string]any
}

// NewConfig creates a new Config.
func NewConfig() Config {
	return Config{}
}

func (c Config) ListFilters() []string {
	var l []string
	for k := range c.filters {
		l = append(l, k)
	}
	sort.Strings(l)
	return l
}
