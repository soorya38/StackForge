package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// StackPRConfig represents the .stackpr.yaml configuration
type StackPRConfig struct {
	Remote   string   `yaml:"remote"`
	Base     string   `yaml:"base"`
	Branches []string `yaml:"branches"`
	SyncMode string   `yaml:"syncMode"` // rebase, merge, reset
}

// InitConfig initializes the configuration
func InitConfig(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".stackpr")
		viper.SetConfigType("yaml")
	}

	viper.SetDefault("remote", "origin")
	viper.SetDefault("syncMode", "rebase")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
			os.Exit(1)
		}
	}
}

// SaveConfig saves the current configuration to .stackpr.yaml
func SaveConfig() error {
	cfg := StackPRConfig{
		Remote:   viper.GetString("remote"),
		Base:     viper.GetString("base"),
		Branches: viper.GetStringSlice("branches"),
		SyncMode: viper.GetString("syncMode"),
	}
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}
	return os.WriteFile(".stackpr.yaml", data, 0644)
}

// Contains checks if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
