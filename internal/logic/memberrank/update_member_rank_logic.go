package memberrank

import (
	"context"

	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/suyuan32/simple-admin-member-rpc/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-member-rpc/types/mms"

	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMemberRankLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMemberRankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMemberRankLogic {
	return &UpdateMemberRankLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMemberRankLogic) UpdateMemberRank(in *mms.MemberRankInfo) (*mms.BaseResp, error) {
	err := l.svcCtx.DB.MemberRank.UpdateOneID(*in.Id).
		SetNotNilName(in.Name).
		SetNotNilCode(in.Code).
		SetNotNilDescription(in.Description).
		SetNotNilRemark(in.Remark).
		Exec(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &mms.BaseResp{Msg: i18n.UpdateSuccess}, nil
}
