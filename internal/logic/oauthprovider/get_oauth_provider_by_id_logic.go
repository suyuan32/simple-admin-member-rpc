package oauthprovider

import (
	"context"
	"github.com/suyuan32/simple-admin-member-rpc/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-member-rpc/types/mms"

	"github.com/suyuan32/simple-admin-common/utils/pointy"

	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetOauthProviderByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOauthProviderByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOauthProviderByIdLogic {
	return &GetOauthProviderByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOauthProviderByIdLogic) GetOauthProviderById(in *mms.IDReq) (*mms.OauthProviderInfo, error) {
	result, err := l.svcCtx.DB.OauthProvider.Get(l.ctx, in.Id)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &mms.OauthProviderInfo{
		Id:           &result.ID,
		CreatedAt:    pointy.GetPointer(result.CreatedAt.UnixMilli()),
		UpdatedAt:    pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		Name:         &result.Name,
		ClientId:     &result.ClientID,
		ClientSecret: &result.ClientSecret,
		RedirectUrl:  &result.RedirectURL,
		Scopes:       &result.Scopes,
		AuthUrl:      &result.AuthURL,
		TokenUrl:     &result.TokenURL,
		AuthStyle:    &result.AuthStyle,
		InfoUrl:      &result.InfoURL,
	}, nil
}
