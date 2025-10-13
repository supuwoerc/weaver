package initialize

import (
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert/yaml"
	"github.com/supuwoerc/weaver/conf"
)

func LoadConfig(viperClient *viper.Viper, consulClient *api.Client) *conf.Config {
	// 映射viper读取到的config到配置对象中
	var globalConfig conf.Config
	if err := viperClient.Unmarshal(&globalConfig); err != nil {
		panic(err)
	}
	// 获取consul中的kv配置来覆盖本地配置
	kv, _, err := consulClient.KV().Get(globalConfig.Consul.KVKey, nil)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(kv.Value, &globalConfig)
	if err != nil {
		panic(err)
	}
	return &globalConfig
}

func LoadConsulConfig(viperClient *viper.Viper) *conf.ConsulConfig {
	// 映射viper读取到的config到配置对象中
	var globalConfig conf.Config
	if err := viperClient.Unmarshal(&globalConfig); err != nil {
		panic(err)
	}
	return &globalConfig.Consul
}
