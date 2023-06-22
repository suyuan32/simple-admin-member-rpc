package member

import (
	"context"
	"github.com/suyuan32/simple-admin-common/utils/pointy"

	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/suyuan32/simple-admin-member-rpc/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-member-rpc/types/mms"

	"github.com/suyuan32/simple-admin-common/utils/uuidx"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMemberByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMemberByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMemberByIdLogic {
	return &GetMemberByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMemberByIdLogic) GetMemberById(in *mms.UUIDReq) (*mms.MemberInfo, error) {
	result, err := l.svcCtx.DB.Member.Get(l.ctx, uuidx.ParseUUIDString(in.Id))
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &mms.MemberInfo{
		Id:        pointy.GetPointer(result.ID.String()),
		CreatedAt: pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt: pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		Status:    pointy.GetPointer(uint32(result.Status)),
		Username:  &result.Username,
		Password:  &result.Password,
		Nickname:  &result.Nickname,
		RankId:    &result.RankID,
		Mobile:    &result.Mobile,
		Email:     &result.Email,
		Avatar:    &result.Avatar,
	}, nil
}
