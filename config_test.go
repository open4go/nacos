package nacos

import (
	"github.com/spf13/viper"
	"os"
	"reflect"
	"testing"
)

// 确保先能够通过curl 命令获取到配置
// curl -X GET "https://nacos.r2day.club/nacos/v1/cs/configs?dataId=auth&group=DEFAULT_GROUP"
func TestGetAuthConfig(t *testing.T) {
	skipUnlessNacosE2E(t)
	tests := []struct {
		name string
		want *viper.Viper
	}{
		{
			"test",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAuthConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAuthConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckGetAuthConfig(t *testing.T) {
	skipUnlessNacosE2E(t)
	tests := []struct {
		name string
		want string
	}{
		{
			"test slice config",
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckGetAuthConfig(); got != tt.want {
				t.Errorf("CheckGetAuthConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 这些测试需要真实可访问的 Nacos 服务端（NACOS_HOST / NACOS_PORT / NACOS_GROUP 等），默认跳过；
// 设置 NACOS_E2E=1 后运行。
func skipUnlessNacosE2E(t *testing.T) {
	t.Helper()
	if os.Getenv("NACOS_E2E") != "1" {
		t.Skip("skip: set NACOS_E2E=1 with a reachable Nacos server to run")
	}
}
