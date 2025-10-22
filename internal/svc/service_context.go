package svc

import (
	"github.com/redis/go-redis/v9"

	"github.com/suyuan32/simple-admin-member-rpc/ent"
	"github.com/suyuan32/simple-admin-member-rpc/internal/config"

	_ "github.com/suyuan32/simple-admin-member-rpc/ent/runtime"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config
	DB     *ent.Client
	Redis  redis.UniversalClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	entOpts := []ent.Option{
		ent.Log(logx.Info),
		ent.Driver(c.DatabaseConf.NewNoCacheDriver()),
	}

	if c.DatabaseConf.Debug {
		entOpts = append(entOpts, ent.Debug())
	}

	db := ent.NewClient(entOpts...)

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  c.RedisConf.MustNewUniversalRedis(),
	}
}
