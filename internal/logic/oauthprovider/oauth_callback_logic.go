package oauthprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/suyuan32/simple-admin-member-rpc/ent/member"
	"github.com/suyuan32/simple-admin-member-rpc/internal/utils/dberrorhandler"
	"github.com/suyuan32/simple-admin-member-rpc/types/mms"

	"github.com/suyuan32/simple-admin-common/msg/logmsg"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/errorx"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/suyuan32/simple-admin-common/i18n"

	"github.com/suyuan32/simple-admin-member-rpc/ent"
	"github.com/suyuan32/simple-admin-member-rpc/ent/oauthprovider"
	"github.com/suyuan32/simple-admin-member-rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type OauthCallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type userInfo struct {
	Email    string `json:"email"`
	NickName string `json:"nickName"`
	Picture  string `json:"picture"`
	Mobile   string `json:"mobile"`
}

func NewOauthCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OauthCallbackLogic {
	return &OauthCallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OauthCallbackLogic) OauthCallback(in *mms.CallbackReq) (*mms.MemberInfo, error) {
	provider := strings.Split(in.State, "-")[1]
	if _, ok := providerConfig[provider]; !ok {
		p, err := l.svcCtx.DB.OauthProvider.Query().Where(oauthprovider.NameEQ(provider)).First(l.ctx)
		if err != nil {
			return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
		}

		providerConfig[provider] = oauth2.Config{
			ClientID:     p.ClientID,
			ClientSecret: p.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   replaceKeywords(p.AuthURL, p),
				TokenURL:  p.TokenURL,
				AuthStyle: oauth2.AuthStyle(p.AuthStyle),
			},
			RedirectURL: p.RedirectURL,
			Scopes:      strings.Split(p.Scopes, " "),
		}
		if _, ok := userInfoURL[p.Name]; !ok {
			userInfoURL[p.Name] = p.InfoURL
		}
	}

	// get user information
	content, err := getUserInfo(providerConfig[provider], userInfoURL[provider], in.Code)
	if err != nil {
		return nil, errorx.NewInvalidArgumentError(err.Error())
	}

	// find or register user
	var u userInfo
	err = json.Unmarshal(content, &u)
	if err != nil {
		return nil, errorx.NewInternalError(err.Error())
	}

	if u.Email != "" {
		result, err := l.svcCtx.DB.Member.Query().Where(member.EmailEQ(u.Email)).WithRanks().First(l.ctx)
		if err != nil {
			switch {
			case ent.IsNotFound(err):
				logx.Errorw(err.Error(), logx.Field("detail", in))
				return nil, errorx.NewInvalidArgumentError("login.userNotExist")
			default:
				logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
				return nil, errorx.NewInternalError(i18n.DatabaseError)
			}
		}

		return &mms.MemberInfo{
			Nickname:  &result.Nickname,
			Avatar:    &result.Avatar,
			Mobile:    &result.Mobile,
			Email:     &result.Email,
			Status:    pointy.GetPointer(uint32(result.Status)),
			Id:        pointy.GetPointer(result.ID.String()),
			Username:  &result.Username,
			RankCode:  &result.Edges.Ranks.Code,
			CreatedAt: pointy.GetPointer(result.CreatedAt.UnixMilli()),
			UpdatedAt: pointy.GetPointer(result.UpdatedAt.UnixMilli()),
		}, nil
	}

	return nil, status.Error(codes.InvalidArgument, i18n.Failed)
}

func getUserInfo(c oauth2.Config, infoURL string, code string) ([]byte, error) {
	var token *oauth2.Token
	var err error

	if strings.Contains(c.Endpoint.AuthURL, "feishu") {
		return GetFeishuUserInfo(c, code)
	} else {
		token, err = c.Exchange(context.Background(), code)
	}
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	var response *http.Response
	if c.Endpoint.AuthStyle == 1 {
		response, err = http.Get(strings.ReplaceAll(infoURL, "TOKEN", token.AccessToken))
		if err != nil {
			return nil, fmt.Errorf("failed to get user's information: %s", err.Error())
		}
	} else if c.Endpoint.AuthStyle == 2 {
		client := &http.Client{}
		request, err := http.NewRequest("GET", infoURL, nil)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", "Bearer "+token.AccessToken)

		response, err = client.Do(request)
		if err != nil {
			return nil, fmt.Errorf("failed getting user info: %s", err.Error())
		}
	}

	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err.Error())
	}

	return contents, nil
}

func replaceKeywords(urlData string, oauthData *ent.OauthProvider) (result string) {
	result = strings.ReplaceAll(urlData, "CLIENT_ID", oauthData.ClientID)
	result = strings.ReplaceAll(result, "SECRET", oauthData.ClientSecret)
	result = strings.ReplaceAll(result, "REDIRECT_URL", oauthData.RedirectURL)
	result = strings.ReplaceAll(result, "SCOPE", oauthData.Scopes)
	return result
}
