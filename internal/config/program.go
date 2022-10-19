package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

const (
	programConfigKeyConfluenceHost      = "confluence_host"
	programConfigKeyConfluenceEmail     = "confluence_email"
	programConfigKeyConfluenceAPIKey    = "confluence_api_key"
	programConfigKeyConfluenceBasicAuth = "confluence_basic_auth"
	programConfigKeyNoConfigFile        = "no_config_file"
)

func LoadOrCreateConfigIfNotExists() error {
	programDir, err := createProgramConfigDirIfNotExists()
	if err != nil {
		return err
	}

	viper.SetDefault(programConfigKeyConfluenceHost, "https://<your-domain>.atlassian.net")
	viper.SetDefault(programConfigKeyConfluenceEmail, "your@email.com")
	viper.SetDefault(programConfigKeyConfluenceAPIKey, "YOUR_API_KEY")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(programDir)

	viper.SetEnvPrefix("LIFECYCLEDOC")
	viper.AutomaticEnv()

	if GetNoConfigFile() == "1" {
		return nil
	}

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

func GetConfluenceBasicAuth() string {
	return viper.GetString(programConfigKeyConfluenceBasicAuth)
}

func GetNoConfigFile() string {
	return viper.GetString(programConfigKeyNoConfigFile)
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
