package main

import (
	"fmt"
	"log"
	"os"

	"github.com/madeiramadeirabr/action-lifecycledoc/internal/config"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema"
	"github.com/madeiramadeirabr/action-lifecycledoc/pkg/schema/parser/yaml"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	if err := config.LoadOrCreateConfigIfNotExists(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "lifecycledoc [lifecycle.yaml file path]",
		Short: "Create lifecycle documentation using lifecycle.yaml file definition",
		Args:  cobra.ExactArgs(1),
		RunE:  process,
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
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

	if err := decoder.Decode(lifecycleFile, schameResolver); err != nil {
		return err
	}

	// @todo: Generate confluence page

	return nil
}
