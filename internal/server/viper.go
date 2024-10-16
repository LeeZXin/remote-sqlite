package server

import (
	"github.com/spf13/viper"
)

var (
	vp *viper.Viper
)

func initViper() {
	vp = viper.New()
	vp.SetConfigType("yaml")
	vp.AddConfigPath(".")
	vp.SetConfigName("application.yaml")
	vp.ReadInConfig()
}
