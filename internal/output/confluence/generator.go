package confluence

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/client/confluence"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/types"
	goconfluence "github.com/virtomize/confluence-go-api"
)

const (
	maxRequestGoroutines = 2
)

type Generator struct {
	client         *confluence.Client
	templateWriter *TemplateWriter
}

type GenerateResult struct {
	Content *goconfluence.Content
	Err     error
}

func NewGenerator(client *confluence.Client, templateRetriver TemplateRetriver) *Generator {
	return &Generator{
		client:         client,
		templateWriter: NewTemplateWriter(templateRetriver),
	}
}

func (g *Generator) Generate(ctx context.Context, schemaResolver schema.Resolver) (<-chan GenerateResult, error) {
	confluence, err := schemaResolver.GetConfluence()
	if err != nil {
		return nil, err
	}

	contentBody, err := g.generateContentBody(ctx, schemaResolver)
	if err != nil {
		return nil, fmt.Errorf("can't generate Confluence Content Body: %w", err)
	}

	var (
		pages = confluence.Pages()

		inputChan  = make(chan *types.ConfluencePage, len(pages))
		resultChan = make(chan GenerateResult, len(pages))

		wg sync.WaitGroup
	)

	defer close(inputChan)

	for i := 0; i < maxRequestGoroutines; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for page := range inputChan {
				g.createOrUpdatePage(ctx, page, contentBody, resultChan)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for i := range pages {
		inputChan <- pages[i]
	}

	return resultChan, nil
}

func (g *Generator) createOrUpdatePage(
	ctx context.Context,
	page *types.ConfluencePage,
	contentBody *goconfluence.Body,
	resultChan chan<- GenerateResult,
) {
	limitedCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	content, err := g.createOrUpdate(limitedCtx, page, contentBody)
	if err != nil {
		resultChan <- GenerateResult{
			Err: fmt.Errorf(
				`can't create or update page "%s" in space "%s" with ancestor "%s": %w`,
				page.Title(),
				page.SpaceKey(),
				page.AncestorID(),
				err,
			),
		}

		return
	}

	if content == nil {
		resultChan <- GenerateResult{
			Err: fmt.Errorf("confluence API did not return a body with the generated documentation information"),
		}

		return
	}

	resultChan <- GenerateResult{
		Content: content,
	}
}

func (g *Generator) generateContentBody(ctx context.Context, schemaResolver schema.Resolver) (*goconfluence.Body, error) {
	body := &strings.Builder{}
	if err := g.templateWriter.Write(body, schemaResolver); err != nil {
		return nil, err
	}

	return &goconfluence.Body{
		Storage: goconfluence.Storage{
			Value:          body.String(),
			Representation: "storage",
		},
	}, nil
}

func (g *Generator) createOrUpdate(
	ctx context.Context,
	page *types.ConfluencePage,
	body *goconfluence.Body,
) (*goconfluence.Content, error) {
	content := g.generateContent(page, body)

	exists, err := g.findAndUpdateCurrentContent(ctx, content, 0)
	if err != nil {
		return nil, err
	}

	if !exists {
		return g.client.CreateContent(ctx, content)
	} else {
		return g.client.UpdateContent(ctx, content)
	}
}

func (g *Generator) generateContent(
	page *types.ConfluencePage,
	body *goconfluence.Body,
) *goconfluence.Content {
	content := &goconfluence.Content{
		Type:  "page",
		Title: page.Title(),
		Space: goconfluence.Space{
			Key: page.SpaceKey(),
		},
		Ancestors: []goconfluence.Ancestor{
			{
				ID: page.AncestorID(),
			},
		},
		Body: *body,
		Version: &goconfluence.Version{
			Number:  1,
			Message: "Documentação gerada (lifecycledoc)",
		},
		// Commented because Confluence was returned:
		// '{"statusCode":500,"message":"java.lang.NullPointerException: null"}'
		// Metadata: &goconfluence.Metadata{
		// 	Properties: &goconfluence.Properties{
		// 		ContentAppearanceDraft: &goconfluence.ContentAppearanceDraft{
		// 			Value: "full-width",
		// 		},
		// 		ContentAppearancePublished: &goconfluence.ContentAppearancePublished{
		// 			Value: "full-width",
		// 		},
		// 	},
		// },
	}

	return content
}

func (g *Generator) findAndUpdateCurrentContent(ctx context.Context, content *goconfluence.Content, start int) (bool, error) {
	result, err := g.client.FindContent(ctx, content, start)
	if err != nil {
		return false, fmt.Errorf("can't search current documentation: %w", err)
	}

	for i := range result.Results {
		if result.Results[i].Title == content.Title {
			for j := range result.Results[i].Ancestors {
				if result.Results[i].Ancestors[j].ID == content.Ancestors[0].ID {
					content.ID = result.Results[i].ID
					content.Version.Number = result.Results[i].Version.Number + 1

					return true, nil
				}
			}
		}
	}

	if result.Size > result.Limit {
		return g.findAndUpdateCurrentContent(ctx, content, start+result.Limit)
	}

	return false, nil
}
