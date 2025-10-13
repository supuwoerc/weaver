package initialize

import (
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/supuwoerc/weaver/conf"
)

func NewConsulClient(conf *conf.ConsulConfig) *api.Client {
	if conf.Address == "" {
		panic("consul address is required")
	}
	cfg := &api.Config{
		Address: conf.Address,
		Scheme:  conf.Scheme,
		Token:   conf.Token,
		TLSConfig: api.TLSConfig{
			CAFile:             conf.CAFile,
			CertFile:           conf.CertFile,
			KeyFile:            conf.KeyFile,
			InsecureSkipVerify: conf.InsecureSkipVerify,
		},
	}
	if conf.WaitTime > 0 {
		cfg.WaitTime = conf.WaitTime * time.Second
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
