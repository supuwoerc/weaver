package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/supuwoerc/weaver/pkg/utils"
)

// ServiceDiscoveryClient 服务发现客户端
type ServiceDiscoveryClient struct {
	discovery    *ServiceDiscovery
	loadBalancer *utils.LoadBalancer[*ServiceInstance]
}

// NewServiceDiscoveryClient 创建服务发现客户端
func NewServiceDiscoveryClient(client *api.Client, logger DiscoveryLogger) *ServiceDiscoveryClient {
	discovery := NewServiceDiscovery(client, logger)
	lb := utils.NewLoadBalancer[*ServiceInstance](discovery, utils.RoundRobin)
	return &ServiceDiscoveryClient{
		discovery:    discovery,
		loadBalancer: lb,
	}
}
