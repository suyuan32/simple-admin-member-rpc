package oauthprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-member-rpc/ent/oauthprovider"
	"github.com/suyuan32/simple-admin-member-rpc/internal/utils/dberrorhandler"
	"github.com/zeromicro/go-zero/core/errorx"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"strings"

	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/suyuan32/simple-admin-member-rpc/types/mms"

	"github.com/zeromicro/go-zero/core/logx"
)

type WechatMiniProgramLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWechatMiniProgramLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatMiniProgramLoginLogic {
	return &WechatMiniProgramLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type WechatMiniProgramLoginResp struct {
	OpenId     string `json:"openid,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
	UnionId    string `json:"unionid,omitempty"`
	ErrorCode  int    `json:"errcode,omitempty"`
	ErrorMsg   string `json:"errmsg,omitempty"`
}

func (l *WechatMiniProgramLoginLogic) WechatMiniProgramLogin(in *mms.OauthLoginReq) (*mms.BaseResp, error) {
	var config oauth2.Config
	if v, ok := providerConfig[in.Provider]; ok {
		config = v
	} else {
		p, err := l.svcCtx.DB.OauthProvider.Query().Where(oauthprovider.NameEQ(in.Provider)).First(l.ctx)
		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}

		providerConfig[p.Name] = oauth2.Config{
			ClientID:     p.ClientID,
			ClientSecret: p.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   p.AuthURL,
				TokenURL:  p.TokenURL,
				AuthStyle: oauth2.AuthStyle(p.AuthStyle),
			},
			RedirectURL: p.RedirectURL,
			Scopes:      strings.Split(p.Scopes, " "),
		}
		config = providerConfig[p.Name]

		if _, ok := userInfoURL[p.Name]; !ok {
			userInfoURL[p.Name] = p.InfoURL
		}

	}

	resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		config.ClientID, config.ClientSecret, in.State))

	if err != nil {
		l.Logger.Errorw("failed to authorize the wechat mini program login request", logx.Field("detail", err))
		return nil, errorx.NewInvalidArgumentError(i18n.Failed)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	var wechatResp WechatMiniProgramLoginResp
	err = json.Unmarshal(body, &wechatResp)
	if err != nil {
		l.Logger.Errorw("failed to unmarshal the response from wechat mini program login resp", logx.Field("detail", err),
			logx.Field("wechatResp", wechatResp))
		return nil, errorx.NewInternalError(i18n.Failed)
	}

	if wechatResp.ErrorCode != 0 {
		l.Logger.Errorw("failed to get the data from wechat, check the code and configuration", logx.Field("wechatResp", wechatResp),
			logx.Field("code", in.State))
		return nil, errorx.NewInvalidArgumentError(i18n.Failed)
	}

	return &mms.BaseResp{Msg: wechatResp.OpenId}, nil

}
