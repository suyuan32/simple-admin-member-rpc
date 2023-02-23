package base

import (
	"context"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/suyuan32/simple-admin-core/pkg/enum"
	"github.com/suyuan32/simple-admin-core/pkg/i18n"
	"github.com/suyuan32/simple-admin-core/pkg/msg/logmsg"
	"github.com/suyuan32/simple-admin-core/pkg/statuserr"
	"github.com/suyuan32/simple-admin-core/pkg/utils"
	"github.com/suyuan32/simple-admin-core/rpc/types/core"
	"github.com/zeromicro/go-zero/core/errorx"

	"github.com/suyuan32/simple-admin-member-rpc/ent"
	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/suyuan32/simple-admin-member-rpc/mms"

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
	err := l.insertApiData()
	if err != nil {
		return nil, statuserr.NewInternalError(err.Error())
	}

	err = l.insertMenuData()
	if err != nil {
		return nil, statuserr.NewInternalError(err.Error())
	}

	if err := l.svcCtx.DB.Schema.Create(l.ctx, schema.WithForeignKeys(false)); err != nil {
		logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
		return nil, errorx.NewCodeError(enum.Internal, err.Error())
	}

	return &mms.BaseResp{
		Msg: i18n.Success,
	}, nil
}

func (l *InitDatabaseLogic) insertApiData() (err error) {

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/member/create",
		Description: "apiDesc.createMember",
		ApiGroup:    "member",
		Method:      "POST",
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/member/update",
		Description: "apiDesc.updateMember",
		ApiGroup:    "member",
		Method:      "POST",
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/member/delete",
		Description: "apiDesc.deleteMember",
		ApiGroup:    "member",
		Method:      "POST",
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/member/list",
		Description: "apiDesc.getMemberList",
		ApiGroup:    "member",
		Method:      "POST",
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/member",
		Description: "apiDesc.getMemberById",
		ApiGroup:    "member",
		Method:      "Post",
	})

	if err != nil {
		return err
	}

	// MEMBER RANK

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/member_rank/create",
		Description: "apiDesc.createMemberRank",
		ApiGroup:    "member_rank",
		Method:      "POST",
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/member_rank/update",
		Description: "apiDesc.updateMemberRank",
		ApiGroup:    "member_rank",
		Method:      "POST",
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/member_rank/delete",
		Description: "apiDesc.deleteMemberRank",
		ApiGroup:    "member_rank",
		Method:      "POST",
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/member_rank/list",
		Description: "apiDesc.getMemberRankList",
		ApiGroup:    "member_rank",
		Method:      "POST",
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateApi(l.ctx, &core.ApiInfo{
		Path:        "/member_rank",
		Description: "apiDesc.getMemberRankById",
		ApiGroup:    "member_rank",
		Method:      "Post",
	})

	if err != nil {
		return err
	}

	return nil
}

func (l *InitDatabaseLogic) insertMenuData() error {
	_, err := l.svcCtx.CoreRpc.CreateMenu(l.ctx, &core.MenuInfo{
		Id:        0,
		CreatedAt: 0,
		UpdatedAt: 0,
		Level:     2,
		ParentId:  16,
		Path:      "",
		Name:      "MemberManagementDirectory",
		Redirect:  "",
		Component: "LAYOUT",
		Sort:      1,
		Disabled:  false,
		Meta: &core.Meta{
			Title:              "route.memberManagement",
			Icon:               "ic:round-person-outline",
			HideMenu:           false,
			HideBreadcrumb:     false,
			IgnoreKeepAlive:    false,
			HideTab:            false,
			FrameSrc:           "",
			CarryParam:         false,
			HideChildrenInMenu: false,
			Affix:              false,
			DynamicLevel:       0,
			RealPath:           "",
		},
		MenuType: 0,
	})

	if err != nil {
		return err
	}

	return err
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
		SetPassword(utils.BcryptEncrypt("simple-admin")),
	)

	members = append(members, l.svcCtx.DB.Member.Create().
		SetUsername("VIPMember").
		SetNickname("VIP Member").
		SetEmail("vip@gmail.com").
		SetMobile("18888888889").
		SetRankID(2).
		SetPassword(utils.BcryptEncrypt("simple-admin")),
	)

	err := l.svcCtx.DB.Member.CreateBulk(members...).Exec(l.ctx)
	if err != nil {
		logx.Errorw(err.Error())
		return statuserr.NewInternalError(err.Error())
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
		return statuserr.NewInternalError(err.Error())
	} else {
		return nil
	}
}
