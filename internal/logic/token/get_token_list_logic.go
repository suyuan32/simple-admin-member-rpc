package token

import (
	"context"
	"github.com/suyuan32/simple-admin-member-rpc/ent/member"
	"github.com/suyuan32/simple-admin-member-rpc/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-member-rpc/types/mms"

	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/suyuan32/simple-admin-common/utils/uuidx"

	"github.com/suyuan32/simple-admin-member-rpc/ent"
	"github.com/suyuan32/simple-admin-member-rpc/ent/predicate"
	"github.com/suyuan32/simple-admin-member-rpc/ent/token"
	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTokenListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenListLogic {
	return &GetTokenListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTokenListLogic) GetTokenList(in *mms.TokenListReq) (*mms.TokenListResp, error) {
	var tokens *ent.TokenPageList
	var err error
	if in.Username == nil && in.Uuid == nil && in.Nickname == nil && in.Email == nil {
		tokens, err = l.svcCtx.DB.Token.Query().Page(l.ctx, in.Page, in.PageSize)

		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}
	} else {
		var predicates []predicate.Member

		if in.Uuid != nil {
			predicates = append(predicates, member.IDEQ(uuidx.ParseUUIDString(*in.Uuid)))
		}

		if in.Username != nil {
			predicates = append(predicates, member.Username(*in.Username))
		}

		if in.Email != nil {
			predicates = append(predicates, member.EmailEQ(*in.Email))
		}

		if in.Nickname != nil {
			predicates = append(predicates, member.NicknameEQ(*in.Nickname))
		}

		u, err := l.svcCtx.DB.Member.Query().Where(predicates...).First(l.ctx)
		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}

		tokens, err = l.svcCtx.DB.Token.Query().Where(token.UUIDEQ(u.ID)).Page(l.ctx, in.Page, in.PageSize)

		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}
	}

	resp := &mms.TokenListResp{}
	resp.Total = tokens.PageDetails.Total

	for _, v := range tokens.List {
		resp.Data = append(resp.Data, &mms.TokenInfo{
			Id:        pointy.GetPointer(v.ID.String()),
			Uuid:      pointy.GetPointer(v.UUID.String()),
			Token:     &v.Token,
			Status:    pointy.GetPointer(uint32(v.Status)),
			Source:    &v.Source,
			Username:  &v.Username,
			ExpiredAt: pointy.GetPointer(v.ExpiredAt.UnixMilli()),
			CreatedAt: pointy.GetPointer(v.CreatedAt.UnixMilli()),
		})
	}

	return resp, nil
}
