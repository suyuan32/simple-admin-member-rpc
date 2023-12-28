package token

import (
	"context"
	"github.com/suyuan32/simple-admin-member-rpc/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-member-rpc/types/mms"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/suyuan32/simple-admin-common/utils/uuidx"

	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTokenByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenByIdLogic {
	return &GetTokenByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTokenByIdLogic) GetTokenById(in *mms.UUIDReq) (*mms.TokenInfo, error) {
	result, err := l.svcCtx.DB.Token.Get(l.ctx, uuidx.ParseUUIDString(in.Id))
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &mms.TokenInfo{
		Id:        pointy.GetPointer(result.ID.String()),
		CreatedAt: pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt: pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		Status:    pointy.GetPointer(uint32(result.Status)),
		Uuid:      pointy.GetPointer(result.UUID.String()),
		Token:     &result.Token,
		Source:    &result.Source,
		Username:  &result.Username,
		ExpiredAt: pointy.GetPointer(result.ExpiredAt.UnixMilli()),
	}, nil
}
