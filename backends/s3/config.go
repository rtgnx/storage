package s3

import (
	"fmt"

	"github.com/spf13/viper"
)

type ConnectionProfile struct {
	URL       string `json:"url"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	API       string `json:"api,omitempty"`
	Path      string `json:"path,omitempty"`
}

type Config struct {
	Version string                       `json:"version,omitempty"`
	Aliases map[string]ConnectionProfile `json:"aliases,omitempty"`
}

var (
	cfg Config
)

func ReadConfig() {
	viper.SetConfigName("config")    // name of config file (without extension)
	viper.SetConfigType("json")      // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/.mc") // call multiple times to add many search paths
	viper.AddConfigPath(".")         // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}
}
