package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jaredwarren/rpi_music/log"
	"github.com/spf13/viper"
)

const (
	ConfigFile = "config"
	ConfigPath = "./config"
)

// InitConfig load config file, write defaults if no file exists.
func InitConfig(logger log.Logger) {
	viper.SetConfigName(ConfigFile) // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			writeDefaultConfig(logger)
		} else {
			logger.Panic("error reading config", log.Error(err))
		}
	}
}

// writeDefaultConfig Set then write config file.
// should only run first time app is launched and no config file is found
func writeDefaultConfig(logger log.Logger) {
	fp := filepath.Join(ConfigPath, fmt.Sprintf("%s.yml", ConfigFile))
	logger.Info("writing default config", log.Any("file_path", fp))
	f, err := os.Create(fp)
	if err != nil {
		logger.Panic("error creating config file", log.Any("file_path", fp), log.Error(err))
	}
	f.Close()

	if err := viper.ReadInConfig(); err != nil {
		logger.Panic("error reading config file", log.Error(err))
	}

	SetDefaults()

	// Write config
	err = viper.WriteConfig()
	if err != nil {
		logger.Panic("error reading config file", log.Error(err))
	}
}

// SetDefaults sets hard-coded default values
func SetDefaults() {
	viper.Set("https", true)
	viper.Set("host", ":8000")
}
