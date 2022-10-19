package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/madeiramadeirabr/action-lifecycledoc/internal/config"
	"github.com/madeiramadeirabr/action-lifecycledoc/internal/output/confluence"
	confluenceRest "github.com/madeiramadeirabr/action-lifecycledoc/pkg/client/confluence"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/parser/yaml"
	"github.com/spf13/cobra"
	goconfluence "github.com/virtomize/confluence-go-api"
)

const (
	titlePrefixFlag  = "titlePrefix"
	outputFormatFlag = "outputFormat"
)

var (
	errLog *log.Logger
)

func init() {
	errLog = log.New(os.Stderr, "", log.Lmicroseconds)

	if err := config.LoadOrCreateConfigIfNotExists(); err != nil {
		errLog.Fatal(err)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "lifecycledoc [lifecycle.yaml file path]",
		Short: "Create lifecycle documentation using lifecycle.yaml file definition",
		Args:  cobra.ExactArgs(1),
		RunE:  process,
	}

	rootCmd.Flags().String(titlePrefixFlag, "", "Specifies a prefix for Confluence page titles")
	rootCmd.Flags().String(outputFormatFlag, "cli", "Specifies the output format. Supported formats: cli, github-action")

	if err := rootCmd.Execute(); err != nil {
		errLog.Fatal(err)
	}
}

func process(cmd *cobra.Command, args []string) error {
	var successWriter successResultWriter

	switch format, _ := cmd.Flags().GetString(outputFormatFlag); format {
	case "cli":
		successWriter = newCLISuccessResultWriter()
	case "github-action":
		successWriter = &githubSuccessResultWriter{}
	default:
		return fmt.Errorf("output format '%s' unknown", format)
	}

	lifecycleFile, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("can't open lifecycle YAML file '%s': %w", args[0], err)
	}

	defer lifecycleFile.Close()

	schameResolver := schema.NewBasicResolver()
	decoder := yaml.NewDecoder()

	titlePrefix, _ := cmd.Flags().GetString(titlePrefixFlag)
	if len(titlePrefix) > 0 {
		schameResolver.SetConfluencePageTitlePrefix(titlePrefix)
	}

	if err := decoder.Decode(lifecycleFile, schameResolver); err != nil {
		return err
	}

	basicAuth := config.GetConfluenceBasicAuth()

	if len(basicAuth) < 1 {
		basicAuth = confluenceRest.GenerateBasicAuthorization(
			config.GetConfluenceEmail(),
			config.GetConfluenceAPIKey(),
		)
	}

	generator := confluence.NewGenerator(
		confluenceRest.NewClient(
			http.DefaultClient,
			config.GetConfluenceHost(),
			basicAuth,
		),
		confluence.NewHTMLTemplateRetriver(),
	)

	resultChan, err := generator.Generate(context.Background(), schameResolver)
	if err != nil {
		return err
	}

	var lastErr error

	for result := range resultChan {
		if result.Err != nil {
			if lastErr != nil {
				lastErr = fmt.Errorf("%s\n\t - %s", lastErr, result.Err)
			} else {
				lastErr = fmt.Errorf("\n\t - %s", result.Err)
			}
		} else {
			successWriter.AddResult(result.Content)
		}
	}

	if lastErr != nil {
		return fmt.Errorf("the following errors occur when creating documentation pages: %s", lastErr)
	}

	if err := successWriter.Output(); err != nil {
		return err
	}

	return nil
}

type successResultWriter interface {
	AddResult(content *goconfluence.Content)
	Output() error
}

type cliSuccessResultWriter struct {
	logger *log.Logger
}

func (c *cliSuccessResultWriter) AddResult(content *goconfluence.Content) {
	c.logger.Printf("documentation generated: %s%s", content.Links.Base, content.Links.TinyUI)
}

func (*cliSuccessResultWriter) Output() error {
	return nil
}

func newCLISuccessResultWriter() *cliSuccessResultWriter {
	return &cliSuccessResultWriter{
		logger: log.New(os.Stdout, "", log.Lmicroseconds),
	}
}

type githubSuccessResultWriter struct {
	links []string
}

func (g *githubSuccessResultWriter) AddResult(content *goconfluence.Content) {
	g.links = append(g.links, fmt.Sprintf("%s%s", content.Links.Base, content.Links.TinyUI))
}

// Output write Confluence page links in GitHub output format
// https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#setting-an-output-parameter
func (g *githubSuccessResultWriter) Output() error {
	fmt.Print("links=")

	if err := json.NewEncoder(os.Stdout).Encode(g.links); err != nil {
		return fmt.Errorf("can't write links json: %w", err)
	}

	return nil
}
