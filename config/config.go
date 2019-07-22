package config

import (
	"github.com/spf13/viper"
	//"os"
)

var conf *viper.Viper = nil

func LoadConfig() *viper.Viper {
	if conf == nil {
		//MODE := os.Getenv("MODE")
		v := viper.New()
		v.AddConfigPath("/conf/")
		v.SetConfigType("yaml")

		v.SetConfigName("dev")
		if err := v.ReadInConfig(); err != nil {
			return nil
		}
		conf = v
	}
	return conf
}
