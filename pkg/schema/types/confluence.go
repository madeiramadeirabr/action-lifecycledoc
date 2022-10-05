package types

import "errors"

type Confluence struct {
	pages []*ConfluencePage
}

func (c *Confluence) Pages() []*ConfluencePage {
	return c.pages
}

func (c *Confluence) AddPage(page *ConfluencePage) {
	c.pages = append(c.pages, page)
}

// ConfluencePage represents a immutable confluence page
type ConfluencePage struct {
	title string
	// spaceKey stores the key of space to publish this page
	spaceKey string
	// acestorID stores the parent page ID
	ancestorID string
}

func (c *ConfluencePage) Title() string {
	return c.title
}

func (c *ConfluencePage) SpaceKey() string {
	return c.spaceKey
}

func (c *ConfluencePage) AncestorID() string {
	return c.ancestorID
}

func NewConfluencePage(title, spaceKey, ancestorID string) (*ConfluencePage, error) {
	if len(title) < 1 {
		return nil, errors.New("title cannot be empty")
	}

	if len(spaceKey) < 1 {
		return nil, errors.New("space key cannot be empty")
	}

	if len(ancestorID) < 1 {
		return nil, errors.New("ancestor id cannot be empty")
	}

	return &ConfluencePage{
		title:      title,
		spaceKey:   spaceKey,
		ancestorID: ancestorID,
	}, nil
}
