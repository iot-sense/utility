package utility

import (
	"bytes"

	"github.com/spf13/viper"
)

//G_CONFIGER global
var G_CONFIGER = setupConfig()

//setupConfig func
func setupConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName("env")
	v.SetConfigType("json")
	v.AddConfigPath("./env/")
	err := v.ReadInConfig()
	if err != nil {
		Logger.Error(err)
	}
	env := v.GetString("env")
	if env == "dev" {
		v.SetConfigName("dev_env")
	} else {
		v.SetConfigName(env + "_env")
	}
	err = v.ReadInConfig()
	if err != nil {
		Logger.Error(err)
	}
	return v
}

//CreateConfig func
func CreateConfig(content []byte) *viper.Viper {
	Logger.Debug("CONTENT: ", string(content))
	v := viper.New()
	v.SetConfigType("json")
	v.ReadConfig(bytes.NewBuffer(content))
	return v
}
