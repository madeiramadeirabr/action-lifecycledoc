package confluence

import "strings"

type outputData struct {
	Types           []*typeOutput
	PublishedEvents []*publishedEventOutput
	ConsumedEvents  []*consumedEventOutput
}

type typeOutput struct {
	Name        string
	Type        string
	Description string
	Nullable    bool
	Format      string
	Enum        []enumValue
	Example     string
}

func (t *typeOutput) ExampleMultipleLine() bool {
	return strings.Contains(t.Example, string('\n'))
}

func (t *typeOutput) TotalAttributes() int {
	// Description + Type
	total := 2

	if t.Nullable {
		total++
	}

	if len(t.Format) > 0 {
		total++
	}

	if len(t.Enum) > 0 {
		total++
	}

	if len(t.Example) > 0 {
		total++
	}

	return total
}

type enumValue struct {
	Value string
	// HasMore indicates if has more items
	HasMore bool
}

type publishedEventOutput struct {
	Name        string
	Visibility  string
	Module      string
	Description string
	Example     string
}

type consumedEventOutput struct {
	Name        string
	Description string
}
