package config

import (
	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigName("configs/config") // name of config file (without extension)
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	err = viper.ReadInConfig()            // Find and read the config file
	return err
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}
