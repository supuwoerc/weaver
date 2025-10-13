package conf

import "time"

type ConsulConfig struct {
	Address            string        `mapstructure:"address"`              // address
	Scheme             string        `mapstructure:"scheme"`               // http or https
	Token              string        `mapstructure:"token"`                // token
	KVKey              string        `mapstructure:"kv_key"`               // kv_key
	WaitTime           time.Duration `mapstructure:"wait_time"`            // wait_time
	CAFile             string        `mapstructure:"ca_file"`              // ca_file (CA 配置（三选一）)
	CertFile           string        `mapstructure:"cert_file"`            // cert_file (CA 配置（三选一）)
	KeyFile            string        `mapstructure:"key_file"`             // key_file  (CA 配置（三选一）)
	InsecureSkipVerify bool          `mapstructure:"insecure_skip_verify"` // insecure_skip_verify
}

//WaitTime: 5 * time.Minute,
//TLSConfig: api.TLSConfig{
//CAFile:             "/etc/ssl/certs/consul-ca.crt",
//CertFile:           "/etc/ssl/certs/consul-client.crt",
//KeyFile:            "/etc/ssl/private/consul-client.key",
//InsecureSkipVerify: false,
//},
