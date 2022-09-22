package types

import "errors"

type Project struct {
	name       string
	confluence *Confluence
}

func (p *Project) Name() string {
	return p.name
}

func (p *Project) Confluence() *Confluence {
	return p.confluence
}

func NewProject(name string) (*Project, error) {
	if len(name) < 1 {
		return nil, errors.New("project name cannot be empty")
	}

	return &Project{
		name:       name,
		confluence: &Confluence{},
	}, nil
}
