package svc

import (
	"github.com/suyuan32/simple-admin-core/rpc/coreclient"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/suyuan32/simple-admin-member-rpc/ent"
	"github.com/suyuan32/simple-admin-member-rpc/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config  config.Config
	DB      *ent.Client
	Redis   *redis.Redis
	CoreRpc coreclient.Core
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := ent.NewClient(
		ent.Log(logx.Info), // logger
		ent.Driver(c.DatabaseConf.NewNoCacheDriver()),
		ent.Debug(), // debug mode
	)

	return &ServiceContext{
		Config:  c,
		DB:      db,
		Redis:   redis.MustNewRedis(c.RedisConf),
		CoreRpc: coreclient.NewCore(zrpc.NewClientIfEnable(c.CoreRpc)),
	}
}
