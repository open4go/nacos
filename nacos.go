package nacos

import (
	"bytes"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"sync"
)

// ConfigHandler 配置句柄
var ConfigHandler config_client.IConfigClient

const (
	defaultConfigType = "yaml"
	defaultScheme     = "http"
	tlsScheme         = "https"
)

type OnChangeFunc []func(conf *viper.Viper)

type hotConfig struct {
	lock             sync.RWMutex
	defaultConfig    map[string]*viper.Viper
	configChangeFunc map[string]OnChangeFunc
}

var (
	defaultConfig = &hotConfig{
		lock:             sync.RWMutex{},
		defaultConfig:    make(map[string]*viper.Viper),
		configChangeFunc: make(map[string]OnChangeFunc),
	}

	// nacos default config
	host             = "localhost"
	port      uint64 = 8848
	namespace        = ""
	group            = ""
	username         = ""
	password         = ""
	grpcPort  uint64 = 9848
)

func RegisterConfigChanged(dataID string, f func(conf *viper.Viper)) {
	defaultConfig.configChangeFunc[dataID] = append(defaultConfig.configChangeFunc[dataID], f)
}

func (c *hotConfig) Read(dataID string) *viper.Viper {
	defer c.lock.RUnlock()
	c.lock.RLock()

	if c.defaultConfig == nil {
		c.defaultConfig = map[string]*viper.Viper{}
	}
	return c.defaultConfig[dataID]
}

func (c *hotConfig) Write(dataID string, data []byte) error {
	c.lock.Lock()
	if c.defaultConfig == nil {
		c.defaultConfig = map[string]*viper.Viper{}
	}

	conf := viper.New()
	conf.SetConfigType(defaultConfigType)
	err := conf.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		c.lock.Unlock()
		return err
	}

	c.defaultConfig[dataID] = conf
	c.lock.Unlock()
	return nil
}

func init() {
	host = os.Getenv("NACOS_HOST")
	portStr := os.Getenv("NACOS_PORT")
	port, _ = strconv.ParseUint(portStr, 10, 64)
	namespace = os.Getenv("NACOS_NAMESPACE")
	group = os.Getenv("NACOS_GROUP")
	username = os.Getenv("NACOS_USERNAME")
	password = os.Getenv("NACOS_PASSWORD")
	grpcPortStr := os.Getenv("NACOS_GRPC_PORT")
	grpcPort, _ = strconv.ParseUint(grpcPortStr, 10, 64)

	err := connectClient()
	if err != nil {
		panic(err)
	}
}

func connectClient() error {
	var err error
	if ConfigHandler != nil {
		return nil
	}

	scheme := defaultScheme
	if port == 443 {
		scheme = tlsScheme
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(
			host,
			port,
			constant.WithScheme(scheme),
			constant.WithContextPath("/nacos"),
		),
	}

	cc := constant.ClientConfig{
		NamespaceId:         namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "error",
		Username:            username,
		Password:            password,
	}

	ConfigHandler, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}
	return nil
}

func Init(namespace string, dataID string) (*viper.Viper, error) {
	logCtx := log.WithField("namespace", namespace).
		WithField("group", group).
		WithField("host", host).
		WithField("port", port).
		WithField("grpcPort", grpcPort).
		WithField("dataID", dataID)

	err := connectClient()
	if err != nil {
		log.WithField("namespace", namespace).WithField("group", group).
			WithField("dataID", dataID).Error(err)
		return nil, err
	}

	content, err := ConfigHandler.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
	})

	if err != nil {
		logCtx.Error(err)
	}

	err = defaultConfig.Write(dataID, []byte(content))
	if err != nil {
		logCtx.Error(err)
	}

	err = ConfigHandler.ListenConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			err := defaultConfig.Write(dataID, []byte(data))
			if err != nil {
				logCtx.Error(err)
			} else {
				logCtx.Info("update success")
			}
		}})

	return defaultConfig.Read(dataID), err
}

func GetConfig(namespace string, dataID string) *viper.Viper {
	var err error
	r := defaultConfig.Read(dataID)
	if r != nil {
		return defaultConfig.Read(dataID)
	} else {
		r, err = Init(namespace, dataID)
		if err != nil {
			log.WithField("namespace", namespace).WithField("group", group).
				WithField("dataID", dataID).Error(err)
			return nil
		}
		return r
	}
}
