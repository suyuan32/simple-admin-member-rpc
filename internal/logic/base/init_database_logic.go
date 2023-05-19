package base

import (
	"context"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/suyuan32/simple-admin-common/enum/errorcode"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-common/msg/logmsg"
	"github.com/suyuan32/simple-admin-common/utils/encrypt"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/suyuan32/simple-admin-member-rpc/ent"
	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/suyuan32/simple-admin-member-rpc/types/mms"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitDatabaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitDatabaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitDatabaseLogic {
	return &InitDatabaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InitDatabaseLogic) InitDatabase(in *mms.Empty) (*mms.BaseResp, error) {

	if err := l.svcCtx.DB.Schema.Create(l.ctx, schema.WithForeignKeys(false)); err != nil {
		logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
		return nil, errorx.NewCodeError(errorcode.Internal, err.Error())
	}

	err := l.insertMemberData()
	if err != nil {
		return nil, errorx.NewInternalError(err.Error())
	}

	err = l.insertMemberRankData()
	if err != nil {
		return nil, errorx.NewInternalError(err.Error())
	}

	return &mms.BaseResp{
		Msg: i18n.Success,
	}, nil
}

// insert init member data
func (l *InitDatabaseLogic) insertMemberData() error {
	var members []*ent.MemberCreate
	members = append(members, l.svcCtx.DB.Member.Create().
		SetUsername("normalMember").
		SetNickname("Normal Member").
		SetEmail("simpleadmin@gmail.com").
		SetMobile("18888888888").
		SetRankID(1).
		SetPassword(encrypt.BcryptEncrypt("simple-admin")),
	)

	members = append(members, l.svcCtx.DB.Member.Create().
		SetUsername("VIPMember").
		SetNickname("VIP Member").
		SetEmail("vip@gmail.com").
		SetMobile("18888888889").
		SetRankID(2).
		SetPassword(encrypt.BcryptEncrypt("simple-admin")),
	)

	err := l.svcCtx.DB.Member.CreateBulk(members...).Exec(l.ctx)
	if err != nil {
		logx.Errorw(err.Error())
		return errorx.NewInternalError(err.Error())
	} else {
		return nil
	}
}

// insert init member rank data
func (l *InitDatabaseLogic) insertMemberRankData() error {
	var memberRanks []*ent.MemberRankCreate
	memberRanks = append(memberRanks, l.svcCtx.DB.MemberRank.Create().
		SetName("memberRank.normal").
		SetCode("001").
		SetDescription("普通会员 | Normal Member").
		SetRemark("普通会员 | Normal Member"),
	)

	memberRanks = append(memberRanks, l.svcCtx.DB.MemberRank.Create().
		SetName("memberRank.vip").
		SetCode("002").
		SetDescription("VIP").
		SetRemark("VIP"),
	)

	err := l.svcCtx.DB.MemberRank.CreateBulk(memberRanks...).Exec(l.ctx)
	if err != nil {
		logx.Errorw(err.Error())
		return errorx.NewInternalError(err.Error())
	} else {
		return nil
	}
}
