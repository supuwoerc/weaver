package initialize

import (
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/supuwoerc/weaver/conf"
)

func NewConsulClient(conf *conf.Config) *api.Client {
	if conf.Consul.Address == "" {
		panic("consul address is required")
	}
	if conf.Consul.Token == "" {
		panic(" No ACL token provided,some operations may fail")
	}
	cfg := &api.Config{
		Address: conf.Consul.Address,
		Scheme:  conf.Consul.Scheme,
		Token:   conf.Consul.Token,
		TLSConfig: api.TLSConfig{
			CAFile:             conf.Consul.CAFile,
			CertFile:           conf.Consul.CertFile,
			KeyFile:            conf.Consul.KeyFile,
			InsecureSkipVerify: conf.Consul.InsecureSkipVerify,
		},
	}
	if conf.Consul.WaitTime > 0 {
		cfg.WaitTime = conf.Consul.WaitTime * time.Second
	}
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// 测试连接
	_, err = client.Agent().Self()
	if err != nil {
		panic(err)
	}
	return client
}
