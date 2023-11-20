package nacos

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const (
	namespaceKey     = "NACOS_NAMESPACE"
	defaultNamespace = "public"

	// 定义dataId
	// 权限验证
	authConf = "auth"
	// 产品服务
	productConf = "product"
)

func GetNamespace() string {
	ns := os.Getenv(namespaceKey)
	if ns != "" {
		return ns
	}
	return defaultNamespace
}

func GetAuthConfig() *viper.Viper {
	return GetConfig(GetNamespace(), authConf)
}

// GetProductConfig 获取产品配置
func GetProductConfig() *viper.Viper {
	return GetConfig(GetNamespace(), productConf)
}

func CheckGetAuthConfig() string {
	c := GetAuthConfig().GetStringMapString("redis")
	for k, v := range c {
		fmt.Println("check slice config item ==> ", k, "val", v)
	}
	return "w"
}
