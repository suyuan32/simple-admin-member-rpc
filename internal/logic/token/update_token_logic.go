package token

import (
	"context"
	"github.com/suyuan32/simple-admin-common/config"
	"github.com/suyuan32/simple-admin-member-rpc/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-member-rpc/types/mms"
	"time"

	"github.com/suyuan32/simple-admin-common/enum/common"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-common/msg/logmsg"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/suyuan32/simple-admin-common/utils/uuidx"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTokenLogic {
	return &UpdateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateTokenLogic) UpdateToken(in *mms.TokenInfo) (*mms.BaseResp, error) {
	token, err := l.svcCtx.DB.Token.UpdateOneID(uuidx.ParseUUIDString(*in.Id)).
		SetNotNilStatus(pointy.GetStatusPointer(in.Status)).
		SetNotNilSource(in.Source).
		Save(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	if uint8(*in.Status) == common.StatusBanned {
		expiredTime := token.ExpiredAt.Sub(time.Now())
		if expiredTime > 0 {
			err = l.svcCtx.Redis.Set(l.ctx, config.RedisTokenPrefix+token.Token, "1", expiredTime).Err()
			if err != nil {
				logx.Errorw(logmsg.RedisError, logx.Field("detail", err.Error()))
				return nil, errorx.NewInternalError(i18n.RedisError)
			}
		}
	} else if uint8(*in.Status) == common.StatusNormal {
		err := l.svcCtx.Redis.Del(l.ctx, config.RedisTokenPrefix+token.Token).Err()
		if err != nil {
			logx.Errorw(logmsg.RedisError, logx.Field("detail", err.Error()))
			return nil, errorx.NewInternalError(i18n.RedisError)
		}
	}

	return &mms.BaseResp{Msg: i18n.UpdateSuccess}, nil
}
