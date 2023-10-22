package nacos

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	publicNamespace = "public"
	authConf        = "auth"
)

func GetAuthConfig() *viper.Viper {
	return GetConfig(publicNamespace, authConf)
}

func CheckGetAuthConfig() string {
	c := GetAuthConfig().GetStringMapString("redis")
	for k, v := range c {
		fmt.Println("check slice config item ==> ", k, "val", v)
	}
	return "w"
}
