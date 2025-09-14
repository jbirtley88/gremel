package util

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig(cfgFile string) error {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("LoadConfig() Could not read file %s: %s", viper.ConfigFileUsed(), err.Error())
	}

	return nil
}
