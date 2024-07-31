package initialize

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"time"
)

func InitCasbin(db *gorm.DB) *casbin.SyncedCachedEnforcer {
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, TablePrefix, "casbin")
	if err != nil {
		panic(err)
	}
	modelString := `
	# 请求规则，r是规则的名称，sub为请求的实体，obj为资源的名称, act为请求的实际动作
	[request_definition]
	r = sub, obj, act
	
	# 权限规则
	[policy_definition]
	p = sub, obj, act
	
	# g 角色的名称，第一个位置为用户，第二个位置为角色，第三个位置为域（在多租户场景下使用）
	[role_definition]
	g = _, _
	
	# 任意一条满足, 就允许访问
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
	`
	m, err := model.NewModelFromString(modelString)
	if err != nil {
		panic(err)
	}
	enforcer, err := casbin.NewSyncedCachedEnforcer(m, adapter)
	if err != nil {
		panic(err)
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		panic(err)
	}
	cacheExpireTime := viper.GetDuration("casbin.cacheExpireTime")
	if cacheExpireTime == 0 {
		cacheExpireTime = 600
	}
	enforcer.SetExpireTime(cacheExpireTime * time.Second)
	return enforcer
}
