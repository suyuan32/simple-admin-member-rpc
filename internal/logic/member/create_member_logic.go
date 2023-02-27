package member

import (
	"context"

	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/suyuan32/simple-admin-member-rpc/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-member-rpc/mms"

	"github.com/suyuan32/simple-admin-core/pkg/i18n"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMemberLogic {
	return &CreateMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMemberLogic) CreateMember(in *mms.MemberInfo) (*mms.BaseUUIDResp, error) {
	result, err := l.svcCtx.DB.Member.Create().
		SetStatus(uint8(in.Status)).
		SetUsername(in.Username).
		SetPassword(in.Password).
		SetNickname(in.Nickname).
		SetRankID(in.RankId).
		SetMobile(in.Mobile).
		SetEmail(in.Email).
		SetAvatar(in.Avatar).
		Save(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(err, in)
	}

	return &mms.BaseUUIDResp{Id: result.ID.String(), Msg: i18n.CreateSuccess}, nil
}
