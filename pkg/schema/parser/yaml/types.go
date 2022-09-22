package yaml

import "gopkg.in/yaml.v2"

type project struct {
	Version string `yaml:"version"`
	Name    string `yaml:"name"`

	Confluence confluence `yaml:"confluence"`

	Events events `yaml:"events"`

	// Types is yaml.MapSlice to keep declaration order
	Types yaml.MapSlice `yaml:"types"`
}

type events struct {
	// Published is yaml.MapSlice to keep declaration order
	Published yaml.MapSlice `yaml:"published"`

	// Consumed is yaml.MapSlice to keep declaration order
	Consumed yaml.MapSlice `yaml:"consumed"`
}

type confluence struct {
	Pages []confluencePage `yaml:"pages"`
}

type confluencePage struct {
	Title      string `yaml:"title"`
	SpaceKey   string `yaml:"spaceKey"`
	AncestorID string `yaml:"ancestorId"`
}
