package base

import (
	"context"
	"github.com/suyuan32/simple-admin-member-rpc/internal/utils/dberrorhandler"

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
		return nil, err
	}

	err = l.insertMemberRankData()
	if err != nil {
		return nil, err
	}

	err = l.insertProviderData()
	if err != nil {
		return nil, err
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
		logx.Errorw("failed to insert member data for initialization", logx.Field("detail", err))
		return dberrorhandler.DefaultEntError(l.Logger, err, nil)
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
		logx.Errorw("failed to insert member rank data for initialization", logx.Field("detail", err))
		return dberrorhandler.DefaultEntError(l.Logger, err, nil)
	} else {
		return nil
	}
}

func (l *InitDatabaseLogic) insertProviderData() error {
	var providers []*ent.OauthProviderCreate

	providers = append(providers, l.svcCtx.DB.OauthProvider.Create().
		SetName("google").
		SetClientID("your client id").
		SetClientSecret("your client secret").
		SetRedirectURL("http://localhost:3100/oauth/login/callback").
		SetScopes("email openid").
		SetAuthURL("https://accounts.google.com/o/oauth2/auth").
		SetTokenURL("https://oauth2.googleapis.com/token").
		SetAuthStyle(1).
		SetInfoURL("https://www.googleapis.com/oauth2/v2/userinfo?access_token=TOKEN"),
	)

	providers = append(providers, l.svcCtx.DB.OauthProvider.Create().
		SetName("github").
		SetClientID("your client id").
		SetClientSecret("your client secret").
		SetRedirectURL("http://localhost:3100/oauth/login/callback").
		SetScopes("email openid").
		SetAuthURL("https://github.com/login/oauth/authorize").
		SetTokenURL("https://github.com/login/oauth/access_token").
		SetAuthStyle(2).
		SetInfoURL("https://api.github.com/user"),
	)

	err := l.svcCtx.DB.OauthProvider.CreateBulk(providers...).Exec(l.ctx)
	if err != nil {
		logx.Errorw("failed to insert member's oauth provider data for initialization", logx.Field("detail", err))
		return dberrorhandler.DefaultEntError(l.Logger, err, nil)
	} else {
		return nil
	}
}
