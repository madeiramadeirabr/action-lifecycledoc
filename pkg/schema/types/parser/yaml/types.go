package yaml

import "gopkg.in/yaml.v2"

type project struct {
	Version string `yaml:"version"`

	// Types is yaml.MapSlice to keep declaration order
	Types yaml.MapSlice `yaml:"types"`
}
