package yaml

import "gopkg.in/yaml.v2"

type project struct {
	Version string `yaml:"version"`

	// Types is yaml.MapSlice to keep declaration order
	Types yaml.MapSlice `yaml:"types"`

	Events events `yaml:"events"`
}

type events struct {
	// Published is yaml.MapSlice to keep declaration order
	Published yaml.MapSlice `yaml:"published"`

	// Consumed is yaml.MapSlice to keep declaration order
	Consumed yaml.MapSlice `yaml:"consumed"`
}
