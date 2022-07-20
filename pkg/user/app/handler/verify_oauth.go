package handler

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"gitlab.haochang.tv/yangyi/examine-code/pkg/user/app/service"
	"gitlab.haochang.tv/yangyi/examine-code/pkg/user/domain"
	"gitlab.haochang.tv/yangyi/examine-code/pkg/utils/oauth"
)

// VerifyOauth 三方登录验证，参数
type VerifyOauth struct {
	Client      oauth.Client `json:"-"`
	RedirectURI string       `json:"redirect_uri" valid:"url,required"`
	// 从三方验证完毕重定向回来时，附带的url query
	RawQuery string     `json:"query" valid:",required"`
	Query    url.Values `json:"-"`
}

// VerifyOauthResult 三方登录验证结果
type VerifyOauthResult struct {
	// 三方验证完毕后，找到的对应系统账号
	// 在没有绑定关系的情况下，会是空值
	User *domain.User
	// 会话凭证
	SessionToken string
	// 三方账号的缓存凭证，后续使用这个凭证可以注册新账号或者绑定已有账号
	VendorToken string
}

// VerifyOauthHandler 三方登录验证
type VerifyOauthHandler struct {
	Oauth      *service.OauthService
	OauthToken *service.OauthTokenService
	Session    *service.SessionTokenService
}

// Handle 验证三方登录
func (h *VerifyOauthHandler) Handle(ctx context.Context, args VerifyOauth) (result VerifyOauthResult, err error) {
	code := args.Query.Get("code")
	if code == "" {
		err = errors.New("empty oauth code")
		return
	}

	vendorUser, err := args.Client.Authorize(code, args.RedirectURI)
	if err != nil {
		err = fmt.Errorf("verify by vendor, %w", err)
		return
	}
	vendorUser.Vendor = args.Client.Vendor()

	user, err := h.Oauth.Find(ctx, vendorUser)
	if errors.Is(err, domain.ErrUserNotFound) {
		var vendorToken string

		vendorToken, err = h.OauthToken.Save(ctx, vendorUser)
		if err != nil {
			err = fmt.Errorf("cache vendor user, %w", err)
			return
		}
		result.VendorToken = vendorToken
		return
	} else if err != nil {
		err = fmt.Errorf("find user by vendor uid, %w", err)
		return
	}

	sessionToken, err := h.Session.Generate(ctx, user)
	if err != nil {
		err = fmt.Errorf("generate session token, %w", err)
		return
	}

	result.User = user
	result.SessionToken = sessionToken
	return
}
