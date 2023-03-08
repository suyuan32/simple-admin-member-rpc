package base

import (
	"context"
	"strings"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/suyuan32/simple-admin-common/enum/common"
	"github.com/suyuan32/simple-admin-common/enum/errorcode"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-common/msg/logmsg"
	"github.com/suyuan32/simple-admin-common/utils/encrypt"
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
		if strings.Contains(err.Error(), "common.createFailed") {
			return nil, errorx.NewInvalidArgumentError(i18n.AlreadyInit)
		}
		return nil, errorx.NewInternalError(err.Error())
	}

	err = l.insertMenuData()
	if err != nil {
		return nil, errorx.NewInternalError(err.Error())
	}

	if err := l.svcCtx.DB.Schema.Create(l.ctx, schema.WithForeignKeys(false)); err != nil {
		logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
		return nil, errorx.NewCodeError(errorcode.Internal, err.Error())
	}

	err = l.insertMemberData()
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
		Method:      "POST",
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
	menuData, err := l.svcCtx.CoreRpc.CreateMenu(l.ctx, &core.MenuInfo{
		Id:        0,
		CreatedAt: 0,
		UpdatedAt: 0,
		Level:     2,
		ParentId:  common.DefaultParentId,
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

	_, err = l.svcCtx.CoreRpc.CreateMenu(l.ctx, &core.MenuInfo{
		Id:        0,
		CreatedAt: 0,
		UpdatedAt: 0,
		Level:     2,
		ParentId:  menuData.Id,
		Path:      "/member",
		Name:      "MemberManagement",
		Redirect:  "",
		Component: "/mms/member/index",
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
		MenuType: 1,
	})

	if err != nil {
		return err
	}

	_, err = l.svcCtx.CoreRpc.CreateMenu(l.ctx, &core.MenuInfo{
		Id:        0,
		CreatedAt: 0,
		UpdatedAt: 0,
		Level:     2,
		ParentId:  menuData.Id,
		Path:      "/member_rank",
		Name:      "MemberRankManagement",
		Redirect:  "",
		Component: "/mms/memberRank/index",
		Sort:      2,
		Disabled:  false,
		Meta: &core.Meta{
			Title:              "route.memberRankManagement",
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
		MenuType: 1,
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
