package service

import (
	"context"
	"database/sql"
	"fmt"
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/initialize"
	"github.com/supuwoerc/weaver/pkg/database"
	"github.com/supuwoerc/weaver/pkg/logger"
	"github.com/supuwoerc/weaver/pkg/utils"

	"gorm.io/gorm"
)

type BasicService struct {
	Logger      *logger.Logger
	DB          *gorm.DB
	Locksmith   *utils.RedisLocksmith
	Conf        *conf.Config
	EmailClient *initialize.EmailClient
}

func NewBasicService(
	logger *logger.Logger,
	db *gorm.DB,
	locksmith *utils.RedisLocksmith,
	conf *conf.Config,
	emailClient *initialize.EmailClient,
) *BasicService {
	return &BasicService{
		Logger:      logger,
		DB:          db,
		Locksmith:   locksmith,
		Conf:        conf,
		EmailClient: emailClient,
	}
}

// Transaction 开启事务,join为true则加入上下文中的事务,如果上下文中没有事务则开启新事务,join为false时直接开启新事务
func (s *BasicService) Transaction(ctx context.Context, join bool, fn database.Action, options ...*sql.TxOptions) error {
	isStarter := false // 是否是事务开启者
	manager := &database.TransactionManager{
		DB:                           s.DB,
		AlreadyCommittedOrRolledBack: false,
	}
	if join {
		if m := database.LoadManager(ctx); m != nil {
			// 从上下文中查找到已经存在的事务
			manager = m
		} else {
			// 未找到已经存在的事务则开启新事务
			isStarter = true
			manager.DB = manager.DB.Begin(options...)
		}
	} else {
		// 开启新事务
		isStarter = true
		manager.DB = manager.DB.Begin(options...)
	}
	if manager.DB.Error != nil {
		return manager.DB.Error
	}
	// 将TransactionManager放入上下文
	wrapContext := database.InjectManager(ctx, manager)
	var execErr error
	wrapFunc := func() {
		defer func() {
			if err := recover(); err != nil {
				stackInfo := string(debug.Stack())
				s.Logger.WithContext(wrapContext).Errorw("transaction recover", "panic", err, "stack", stackInfo)
				execErr = fmt.Errorf("transaction panic: %s", err)
			}
		}()
		execErr = fn(wrapContext)
	}
	wrapFunc()
	if execErr != nil {
		if !manager.AlreadyCommittedOrRolledBack {
			manager.AlreadyCommittedOrRolledBack = true
			if rollback := manager.DB.Rollback(); rollback.Error != nil {
				s.Logger.WithContext(wrapContext).Errorw("rollback fail",
					"err", rollback.Error.Error(),
					"execErr", execErr.Error(),
				)
				return errors.WithMessage(rollback.Error, execErr.Error())
			}
		}
		return execErr
	}
	if isStarter && !manager.AlreadyCommittedOrRolledBack {
		manager.AlreadyCommittedOrRolledBack = true
		if commit := manager.DB.Commit(); commit.Error != nil {
			return commit.Error
		}
	}
	return nil
}
