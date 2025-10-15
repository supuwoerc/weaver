package consul

import (
	"context"
	"fmt"

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

type ServiceCallFunc func(ctx context.Context, ins *ServiceInstance) (any, error)

func (sd *ServiceDiscoveryClient) CallService(ctx context.Context, serviceName string, fn ServiceCallFunc) (any, error) {
	instance, err := sd.loadBalancer.GetServiceInstance(serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get service instance: %w", err)
	}
	defer sd.loadBalancer.ReleaseConnection(instance.ID)
	return fn(ctx, instance)
}
