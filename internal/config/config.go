package config

import (
	"github.com/suyuan32/simple-admin-core/pkg/config"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DatabaseConf config.DatabaseConf
	RedisConf    redis.RedisConf
	CoreRpc      zrpc.RpcClientConf
}
