package main

import (
	"context"
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
)

const (
	titlePrefixFlag = "titlePrefix"
)

var (
	infoLog *log.Logger
	errLog  *log.Logger
)

func init() {
	infoLog = log.New(os.Stdout, "", log.Lmicroseconds)
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

	if err := rootCmd.Execute(); err != nil {
		errLog.Fatal(err)
	}
}

func process(cmd *cobra.Command, args []string) error {
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

	generator := confluence.NewGenerator(
		confluenceRest.NewClient(
			http.DefaultClient,
			config.GetConfluenceHost(),
			config.GetConfluenceEmail(),
			config.GetConfluenceAPIKey(),
		),
		confluence.NewHTMLTemplateRetriver(),
	)

	resultChan, err := generator.Generate(context.Background(), schameResolver)
	if err != nil {
		return err
	}

	for result := range resultChan {
		if result.Err != nil {
			errLog.Print(result.Err)
		} else {
			infoLog.Printf("documentation generated: %s%s", result.Content.Links.Base, result.Content.Links.TinyUI)
		}
	}

	return nil
}
