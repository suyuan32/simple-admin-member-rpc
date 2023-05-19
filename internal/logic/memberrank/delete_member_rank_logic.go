package memberrank

import (
	"context"

	"github.com/suyuan32/simple-admin-member-rpc/ent/memberrank"
	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/suyuan32/simple-admin-member-rpc/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-member-rpc/types/mms"

	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMemberRankLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMemberRankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMemberRankLogic {
	return &DeleteMemberRankLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMemberRankLogic) DeleteMemberRank(in *mms.IDsReq) (*mms.BaseResp, error) {
	_, err := l.svcCtx.DB.MemberRank.Delete().Where(memberrank.IDIn(in.Ids...)).Exec(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &mms.BaseResp{Msg: i18n.DeleteSuccess}, nil
}
