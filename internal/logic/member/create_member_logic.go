package member

import (
	"context"

	"github.com/suyuan32/simple-admin-common/utils/encrypt"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/suyuan32/simple-admin-member-rpc/ent/member"
	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/suyuan32/simple-admin-member-rpc/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-member-rpc/types/mms"

	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/zeromicro/go-zero/core/errorx"
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
	if in.Mobile != nil {
		checkMobile, err := l.svcCtx.DB.Member.Query().Where(member.MobileEQ(*in.Mobile)).Exist(l.ctx)
		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}

		if checkMobile {
			return nil, errorx.NewInvalidArgumentError("login.mobileExist")
		}
	}

	if in.Email != nil {
		checkEmail, err := l.svcCtx.DB.Member.Query().Where(member.EmailEQ(*in.Email)).Exist(l.ctx)
		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}

		if checkEmail {
			return nil, errorx.NewInvalidArgumentError("login.signupUserExist")
		}
	}

	query := l.svcCtx.DB.Member.Create().
		SetNotNilStatus(pointy.GetStatusPointer(in.Status)).
		SetNotNilUsername(in.Username).
		SetNotNilNickname(in.Nickname).
		SetNotNilRankID(in.RankId).
		SetNotNilMobile(in.Mobile).
		SetNotNilEmail(in.Email).
		SetNotNilAvatar(in.Avatar).
		SetNotNilWechatOpenID(in.WechatId).
		SetNotNilExpiredAt(pointy.GetTimeMilliPointer(in.ExpiredAt))

	if in.Password != nil {
		query.SetNotNilPassword(pointy.GetPointer(encrypt.BcryptEncrypt(*in.Password)))
	}
	result, err := query.Save(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &mms.BaseUUIDResp{Id: result.ID.String(), Msg: i18n.CreateSuccess}, nil
}
