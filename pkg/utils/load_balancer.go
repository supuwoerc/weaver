package utils

import (
	"fmt"
	"math/rand"
	"sync"
)

// LoadBalanceStrategy 负载均衡策略
type LoadBalanceStrategy int

const (
	RoundRobin       LoadBalanceStrategy = iota // 轮询
	Random                                      // 随机
	LeastConnections                            // 最少连接
)

type Discovery[T any] interface {
	DiscoverServices(s string) ([]T, error)
}

type Instance interface {
	InstanceID() string
	Address() string
}

// LoadBalancer 负载均衡器
type LoadBalancer[T Instance] struct {
	discovery    Discovery[T]
	strategy     LoadBalanceStrategy
	currentIndex map[string]int // 轮询索引位置
	connCounts   map[string]int // 连接数统计
	mutex        sync.RWMutex
}

// NewLoadBalancer 创建负载均衡器
func NewLoadBalancer[T Instance](discovery Discovery[T], strategy LoadBalanceStrategy) *LoadBalancer[T] {
	return &LoadBalancer[T]{
		discovery:    discovery,
		strategy:     strategy,
		currentIndex: make(map[string]int),
		connCounts:   make(map[string]int),
	}
}

// GetServiceInstance 获取服务实例
func (lb *LoadBalancer[T]) GetServiceInstance(serviceName string) (T, error) {
	var zero T
	instances, err := lb.discovery.DiscoverServices(serviceName)
	if err != nil {
		return zero, err
	}
	if len(instances) == 0 {
		return zero, fmt.Errorf("no healthy instances available for service: %s", serviceName)
	}
	switch lb.strategy {
	case RoundRobin:
		return lb.roundRobin(serviceName, instances)
	case Random:
		return lb.random(instances)
	case LeastConnections:
		return lb.leastConnections(instances)
	default:
		return lb.roundRobin(serviceName, instances)
	}
}

// roundRobin 轮询策略
func (lb *LoadBalancer[T]) roundRobin(serviceName string, instances []T) (T, error) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	index := lb.currentIndex[serviceName]
	instance := instances[index]
	// 更新索引
	lb.currentIndex[serviceName] = (index + 1) % len(instances)
	lb.connCounts[instance.InstanceID()]++
	return instance, nil
}

// random 随机策略
func (lb *LoadBalancer[T]) random(instances []T) (T, error) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	index := rand.Intn(len(instances))
	instance := instances[index]
	lb.connCounts[instance.InstanceID()]++
	return instance, nil
}

// leastConnections 最少连接策略
func (lb *LoadBalancer[T]) leastConnections(instances []T) (T, error) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	var selectedInstance T
	minConnections := -1
	ok := false
	for _, instance := range instances {
		connCount := lb.connCounts[instance.InstanceID()]
		if minConnections == -1 || connCount < minConnections {
			minConnections = connCount
			selectedInstance = instance
			ok = true
		}
	}
	if ok {
		lb.connCounts[selectedInstance.InstanceID()]++
	}
	return selectedInstance, nil
}

// ReleaseConnection 释放连接
func (lb *LoadBalancer[T]) ReleaseConnection(instanceID string) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	if count, exists := lb.connCounts[instanceID]; exists && count > 0 {
		lb.connCounts[instanceID] = count - 1
	}
}

// GetInstanceStats 获取统计信息
func (lb *LoadBalancer[T]) GetInstanceStats(serviceName string) (map[string]int, error) {
	instances, err := lb.discovery.DiscoverServices(serviceName)
	if err != nil {
		return nil, err
	}
	stats := make(map[string]int)
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()
	for _, instance := range instances {
		stats[instance.Address()] = lb.connCounts[instance.InstanceID()]
	}
	return stats, nil
}
