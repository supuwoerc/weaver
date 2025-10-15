package consul

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/supuwoerc/weaver/pkg/logger"
)

// ServiceInstance 服务实例
type ServiceInstance struct {
	ID      string
	Name    string
	Address string
	Port    int
	Tags    []string
	Meta    map[string]string
}

func (s *ServiceInstance) GetInstanceID() string {
	return s.ID
}

func (s *ServiceInstance) GetAddress() string {
	return s.Address
}

type DiscoveryLogger interface {
	logger.LogCtxInterface
}

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	client     *api.Client
	cache      map[string][]*ServiceInstance
	cacheTTL   time.Duration
	lastUpdate time.Time
	logger     DiscoveryLogger
	mutex      sync.RWMutex
}

// NewServiceDiscovery 创建服务发现实例
func NewServiceDiscovery(client *api.Client, logger DiscoveryLogger) *ServiceDiscovery {
	return &ServiceDiscovery{
		client:   client,
		cache:    make(map[string][]*ServiceInstance),
		cacheTTL: 30 * time.Second,
		logger:   logger,
	}
}

// DiscoverServices 发现健康的服务实例
func (sd *ServiceDiscovery) DiscoverServices(serviceName string) ([]*ServiceInstance, error) {
	// 检查缓存
	if instances, ok := sd.getFromCache(serviceName); ok {
		return instances, nil
	}
	// 从 Consul 获取
	instances, err := sd.fetchHealthyServices(serviceName)
	if err != nil {
		return nil, err
	}
	// 更新缓存
	sd.updateCache(serviceName, instances)
	return instances, nil
}

// fetchHealthyServices 从 Consul 获取健康服务
func (sd *ServiceDiscovery) fetchHealthyServices(serviceName string) ([]*ServiceInstance, error) {
	// 使用健康检查 API 获取健康实例
	entries, _, err := sd.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query healthy services: %w", err)
	}
	instances := make([]*ServiceInstance, 0, len(entries))
	for _, entry := range entries {
		service := entry.Service
		instance := &ServiceInstance{
			ID:      service.ID,
			Name:    service.Service,
			Address: service.Address,
			Port:    service.Port,
			Tags:    service.Tags,
			Meta:    service.Meta,
		}
		instances = append(instances, instance)
	}
	sd.logger.WithContext(context.Background()).Infow("Discovered  healthy instances for service",
		"serviceName", serviceName, "count", len(instances))
	return instances, nil
}

// 缓存相关方法
func (sd *ServiceDiscovery) getFromCache(serviceName string) ([]*ServiceInstance, bool) {
	sd.mutex.RLock()
	defer sd.mutex.RUnlock()
	// 检查缓存是否过期
	if time.Since(sd.lastUpdate) > sd.cacheTTL {
		return nil, false
	}
	instances, exists := sd.cache[serviceName]
	return instances, exists && len(instances) > 0
}
func (sd *ServiceDiscovery) updateCache(serviceName string, instances []*ServiceInstance) {
	sd.mutex.Lock()
	defer sd.mutex.Unlock()
	sd.cache[serviceName] = instances
	sd.lastUpdate = time.Now()
}

// ClearCache 清除缓存
func (sd *ServiceDiscovery) ClearCache() {
	sd.mutex.Lock()
	defer sd.mutex.Unlock()
	sd.cache = make(map[string][]*ServiceInstance)
	sd.lastUpdate = time.Time{}
}
