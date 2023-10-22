package nacos

import (
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

// 确保先能够通过curl 命令获取到配置
// curl -X GET "https://nacos.r2day.club/nacos/v1/cs/configs?dataId=auth&group=DEFAULT_GROUP"
func TestGetAuthConfig(t *testing.T) {
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
