package cmd

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	config     Config
	params     Parameters
	cmdParams  Parameters
	rootCmd    = &cobra.Command{
		Use:   "bmctl",
		Short: "bmctl is a command line tool for interacting with BrandMeister API",
		Long:  `bmctl is a command line tool for interacting with BrandMeister API`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.config/bmctl/config.yaml)")
	rootCmd.PersistentFlags().IntVarP(&cmdParams.DeviceID, "device-id", "d", 0, "device ID")
	rootCmd.PersistentFlags().StringVarP(&cmdParams.APIKey, "api-key", "k", "", "API key")

	rootCmd.AddCommand(versionCmd)
}

// Config represents the structure of the entire configuration file.
type Config struct {
	Version       string                   `mapstructure:"version"`
	APIKeys       map[string]string        `mapstructure:"apiKeys"`
	Devices       map[string]*DeviceConfig `mapstructure:"devices"`
	DefaultDevice string                   `mapstructure:"defaultDevice"`
}

// DeviceConfig represents the configuration for a single device.
type DeviceConfig struct {
	DeviceID  int    `mapstructure:"deviceID"`
	APIKey    string `mapstructure:"apiKey"`
	APIKeyRef string `mapstructure:"apiKeyRef"`
}

// Parameters represents the parameters for a single device.
type Parameters struct {
	DeviceID int
	APIKey   string
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)

		k := koanf.New(".")
		configFile = path.Join(homeDir, ".config", "bmctl", "config.yaml")
		err = k.Load(file.Provider(configFile), yaml.Parser())
		if err != nil {
			err = k.Unmarshal("", &config)
			cobra.CheckErr(err)

			for deviceName, deviceConfig := range config.Devices {
				if deviceConfig.APIKeyRef != "" {
					apiKey, ok := config.APIKeys[deviceConfig.APIKeyRef]
					if !ok {
						cobra.CheckErr(fmt.Errorf("API key reference %s for device %s not found", deviceConfig.APIKeyRef, deviceName))
					}
					deviceConfig.APIKey = apiKey
				}
			}

			if config.DefaultDevice != "" {
				defaultDevice, ok := config.Devices[config.DefaultDevice]
				if !ok {
					cobra.CheckErr(fmt.Errorf("default device %s not found", config.DefaultDevice))
				}
				params.DeviceID = defaultDevice.DeviceID
				params.APIKey = defaultDevice.APIKey
			}
		}

		if deviceID := os.Getenv("BMCTL_DEVICE_ID"); deviceID != "" {
			value, err := strconv.ParseInt(deviceID, 10, 64)
			if err != nil {
				cobra.CheckErr(fmt.Errorf("invalid device ID %s set using BMCTL_DEVICE_ID: %v", deviceID, err))
			}
			params.DeviceID = int(value)
		}

		if apiKey := os.Getenv("BMCTL_API_KEY"); apiKey != "" {
			params.APIKey = apiKey
		}

		if cmdParams.DeviceID != 0 {
			params.DeviceID = cmdParams.DeviceID
		}

		if cmdParams.APIKey != "" {
			params.APIKey = cmdParams.APIKey
		}
	}
}
