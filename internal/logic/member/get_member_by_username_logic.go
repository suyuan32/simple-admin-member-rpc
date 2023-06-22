package member

import (
	"context"
	"github.com/suyuan32/simple-admin-common/utils/pointy"

	"github.com/suyuan32/simple-admin-member-rpc/ent/member"
	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/suyuan32/simple-admin-member-rpc/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-member-rpc/types/mms"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMemberByUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMemberByUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMemberByUsernameLogic {
	return &GetMemberByUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMemberByUsernameLogic) GetMemberByUsername(in *mms.UsernameReq) (*mms.MemberInfo, error) {
	result, err := l.svcCtx.DB.Member.Query().Where(member.UsernameEQ(in.Username)).First(l.ctx)
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
