package member

import (
	"context"

	"github.com/suyuan32/simple-admin-member-rpc/ent/member"
	"github.com/suyuan32/simple-admin-member-rpc/ent/predicate"
	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/suyuan32/simple-admin-member-rpc/mms"

	"github.com/suyuan32/simple-admin-core/pkg/i18n"
	"github.com/suyuan32/simple-admin-core/pkg/statuserr"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMemberListLogic {
	return &GetMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMemberListLogic) GetMemberList(in *mms.MemberListReq) (*mms.MemberListResp, error) {
	var predicates []predicate.Member
	if in.Username != "" {
		predicates = append(predicates, member.UsernameContains(in.Username))
	}
	if in.Nickname != "" {
		predicates = append(predicates, member.NicknameContains(in.Nickname))
	}
	if in.Mobile != "" {
		predicates = append(predicates, member.MobileContains(in.Mobile))
	}
	if in.Email != "" {
		predicates = append(predicates, member.EmailContains(in.Email))
	}
	if in.RankId != 0 {
		predicates = append(predicates, member.RankIDEQ(in.RankId))
	}

	result, err := l.svcCtx.DB.Member.Query().Where(predicates...).Page(l.ctx, in.Page, in.PageSize)

	if err != nil {
		logx.Error(err.Error())
		return nil, statuserr.NewInternalError(i18n.DatabaseError)
	}

	resp := &mms.MemberListResp{}
	resp.Total = result.PageDetails.Total

	for _, v := range result.List {
		resp.Data = append(resp.Data, &mms.MemberInfo{
			Id:        v.ID.String(),
			CreatedAt: v.CreatedAt.UnixMilli(),
			UpdatedAt: v.UpdatedAt.UnixMilli(),
			Status:    uint32(v.Status),
			Username:  v.Username,
			Password:  v.Password,
			Nickname:  v.Nickname,
			RankId:    v.RankID,
			Mobile:    v.Mobile,
			Email:     v.Email,
			Avatar:    v.Avatar,
		})
	}

	return resp, nil
}
