package service

import (
	"context"
	"errors"
	"gin-web/pkg/constant"
	"gin-web/pkg/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"runtime/debug"
	"sync"
)

type BasicService struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

var (
	basicOnce sync.Once
	basic     *BasicService
)

func NewBasicService() *BasicService {
	basicOnce.Do(func() {
		basic = &BasicService{
			logger: global.Logger,
			db:     global.DB,
		}
	})
	return basic
}

type action func(ctx context.Context) error

// Start 开启一个新的事务
func (s *BasicService) Start(ctx context.Context, fn action) error {
	tx := s.db.Begin().WithContext(ctx)
	defer func() {
		if err := recover(); err != nil {
			s.logger.Errorf("Transaction panic,堆栈信息:", string(debug.Stack()))
			tx.Rollback()
		}
	}()
	wrapContext := context.WithValue(ctx, constant.TransactionKey, tx)
	if err := fn(wrapContext); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// Join 加入到上下文中的事务
func (s *BasicService) Join(ctx context.Context, fn action) error {
	value := ctx.Value(constant.TransactionKey)
	if value == nil {
		return s.Start(ctx, fn)
	}
	tx, ok := value.(*gorm.DB)
	if !ok {
		return errors.New("获取到上下文中的事务类型不属于gorm.DB")
	}
	defer func() {
		if err := recover(); err != nil {
			s.logger.Errorf("Transaction Join panic,堆栈信息:", string(debug.Stack()))
			tx.Rollback()
		}
	}()
	if err := fn(ctx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
