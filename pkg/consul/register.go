package consul

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/supuwoerc/weaver/pkg/constant"
	"github.com/supuwoerc/weaver/pkg/logger"
)

type RegisterEmailClient interface {
	Alarm2Admin(ctx context.Context, subject constant.Subject, body string) error
}

type RegisterLogger interface {
	logger.LogCtxInterface
}

// ServiceRegister 服务注册器
type ServiceRegister struct {
	client      *api.Client
	serviceID   string
	registered  bool
	logger      RegisterLogger
	emailClient RegisterEmailClient
	mutex       sync.RWMutex
}

// ServiceOption 服务注册选项
type ServiceOption func(*api.AgentServiceRegistration)

// NewServiceRegistry 创建服务注册器
func NewServiceRegistry(client *api.Client, emailClient RegisterEmailClient, logger RegisterLogger) *ServiceRegister {
	return &ServiceRegister{
		client:      client,
		registered:  false,
		emailClient: emailClient,
		logger:      logger,
	}
}

// Register 注册服务
func (sr *ServiceRegister) Register(serviceName, address, protocol string, port int, opts ...ServiceOption) error {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	if sr.registered {
		return fmt.Errorf("service already registered")
	}
	// 生成服务ID
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	serviceID := fmt.Sprintf("%s-%s-%d", serviceName, hostname, port)
	// 创建服务注册配置
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Tags:    []string{},
		Meta: map[string]string{
			"hostname": hostname,
			"started":  time.Now().Format(time.DateTime),
			"protocol": protocol,
		},
	}
	// 应用选项
	for _, opt := range opts {
		opt(registration)
	}
	// 设置默认健康检查（如果未设置）
	if registration.Check == nil {
		registration.Check = &api.AgentServiceCheck{
			TCP:                            fmt.Sprintf("%s:%d", address, port),
			Timeout:                        "3s",
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",
		}
	}
	// 注册服务
	if err = sr.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}
	sr.serviceID = serviceID
	sr.registered = true
	sr.logger.WithContext(context.Background()).Infow("Service registered successfully",
		"service_name", serviceName, "service_id", serviceID)
	return nil
}

// Deregister 注销服务
func (sr *ServiceRegister) Deregister() error {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	if !sr.registered {
		return nil
	}
	if err := sr.client.Agent().ServiceDeregister(sr.serviceID); err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}
	sr.registered = false
	sr.logger.WithContext(context.Background()).Infow("Service deregistered",
		"service_id", sr.serviceID)
	return nil
}

// KeepAlive 保持服务活跃状态
func (sr *ServiceRegister) KeepAlive(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return sr.Deregister()
		case <-ticker.C:
			// 发送心跳或更新 TTL 检查
			if err := sr.updateTTL(); err != nil {
				sr.logger.WithContext(ctx).Errorw("Failed to update TTL, Try to register again", "err", err)
				// 尝试重新注册
				if err = sr.reregister(); err != nil {
					sr.logger.WithContext(ctx).Errorw("Failed to reregister", "err", err)
					if alarmErr := sr.emailClient.Alarm2Admin(ctx, constant.ServiceRegister, err.Error()); alarmErr != nil {
						sr.logger.WithContext(ctx).Errorw("Failed to alarm to admin", "err", alarmErr)
					}
				}
			}
		}
	}
}

// 对于 TTL 检查，更新状态
func (sr *ServiceRegister) updateTTL() error {
	return sr.client.Agent().UpdateTTL(sr.serviceID, "service is healthy", "pass")
}

func (sr *ServiceRegister) reregister() error {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	if !sr.registered {
		return fmt.Errorf("service not registered")
	}
	// 获取当前服务信息并重新注册
	services, err := sr.client.Agent().Services()
	if err != nil {
		return err
	}
	service, exists := services[sr.serviceID]
	if !exists {
		return fmt.Errorf("service not found in consul")
	}
	registration := &api.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Service,
		Address: service.Address,
		Port:    service.Port,
		Tags:    service.Tags,
		Meta:    service.Meta,
	}
	return sr.client.Agent().ServiceRegister(registration)
}

// WithTags 设置服务标签
func WithTags(tags ...string) ServiceOption {
	return func(reg *api.AgentServiceRegistration) {
		reg.Tags = append(reg.Tags, tags...)
	}
}

// WithMeta 设置服务元数据
func WithMeta(meta map[string]string) ServiceOption {
	return func(reg *api.AgentServiceRegistration) {
		if reg.Meta == nil {
			reg.Meta = make(map[string]string)
		}
		for k, v := range meta {
			reg.Meta[k] = v
		}
	}
}

// WithHTTPCheck 设置 HTTP 健康检查
func WithHTTPCheck(scheme, path string, interval time.Duration) ServiceOption {
	return func(reg *api.AgentServiceRegistration) {
		reg.Check = &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("%s://%s:%d%s", scheme, reg.Address, reg.Port, path),
			Timeout:                        "3s",
			Interval:                       interval.String(),
			DeregisterCriticalServiceAfter: "1m",
		}
	}
}

// WithTCPCheck 设置 TCP 健康检查
func WithTCPCheck(interval time.Duration) ServiceOption {
	return func(reg *api.AgentServiceRegistration) {
		reg.Check = &api.AgentServiceCheck{
			TCP:                            fmt.Sprintf("%s:%d", reg.Address, reg.Port),
			Timeout:                        "3s",
			Interval:                       interval.String(),
			DeregisterCriticalServiceAfter: "1m",
		}
	}
}

// WithTTLCheck 设置 TTL 健康检查
func WithTTLCheck(interval time.Duration) ServiceOption {
	return func(reg *api.AgentServiceRegistration) {
		reg.Check = &api.AgentServiceCheck{
			TTL:                            interval.String(),
			DeregisterCriticalServiceAfter: "1m",
		}
	}
}

// WithGRPCCheck 设置 gRPC 健康检查
func WithGRPCCheck(interval time.Duration) ServiceOption {
	return func(reg *api.AgentServiceRegistration) {
		reg.Check = &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", reg.Address, reg.Port),
			Interval:                       interval.String(),
			DeregisterCriticalServiceAfter: "1m",
		}
	}
}
