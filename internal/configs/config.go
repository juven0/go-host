package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Resource struct {
	Name            string
	Endpoint        string
	Destination_URL string
}

type Configuration struct {
	Server struct {
		Host        string
		Listen_port string
	}
	Resources []Resource
}

func NewConfiguration() (*Configuration, error) {
	var Config *Configuration
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %s", err)
	}

	// Remonter d'un niveau si nous sommes dans le dossier 'cmd'
	if filepath.Base(workDir) == "cmd" {
		workDir = filepath.Dir(workDir)
	}

	// Définir le chemin du dossier 'data'
	configPath := filepath.Join(workDir, "data")
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config file: %s", err)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}

	return Config, nil
}
