package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

const (
	programConfigKeyConfluenceHost   = "confluence_host"
	programConfigKeyConfluenceEmail  = "confluence_email"
	programConfigKeyConfluenceAPIKey = "confluence_API_key"

	programConfigDefaultConfluenceHost   = "https://<your-domain>.atlassian.net"
	programConfigDefaultConfluenceEmail  = "your@email.com"
	programConfigDefaultConfluenceAPIKey = "YOUR_API_KEY"
)

func LoadOrCreateConfigIfNotExists() error {
	programDir, err := createProgramConfigDirIfNotExists()
	if err != nil {
		return err
	}

	viper.SetDefault(programConfigKeyConfluenceHost, programConfigDefaultConfluenceHost)
	viper.SetDefault(programConfigKeyConfluenceEmail, programConfigDefaultConfluenceEmail)
	viper.SetDefault(programConfigKeyConfluenceAPIKey, programConfigDefaultConfluenceAPIKey)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(programDir)

	viper.SetEnvPrefix("LIFECYCLEDOC")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, is := err.(viper.ConfigFileNotFoundError); is {
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("can't write default configs to program configuration file: %w", err)
			}

			return fmt.Errorf("please, update your config file: %s/config.yaml", programDir)
		}

		return fmt.Errorf("unable to read settings: %w", err)
	}

	return nil
}

func GetConfluenceHost() string {
	return viper.GetString(programConfigKeyConfluenceHost)
}

func GetConfluenceEmail() string {
	return viper.GetString(programConfigKeyConfluenceEmail)
}

func GetConfluenceAPIKey() string {
	return viper.GetString(programConfigKeyConfluenceAPIKey)
}

func createProgramConfigDirIfNotExists() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("can't retrive user home directory: %w", err)
	}

	programDir := path.Join(homeDir, ".lifecycledoc")
	if _, err := os.Stat(programDir); err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("can't check if directory '%s' exists: %w", programDir, err)
		}

		if err := os.Mkdir(programDir, 0700); err != nil {
			return "", fmt.Errorf("can't create directory '%s': %w", programDir, err)
		}
	}

	return programDir, nil
}
