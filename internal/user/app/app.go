package app

import (
	"ddd-example/internal/option"
	"ddd-example/internal/user/app/adapter"
	"ddd-example/internal/user/app/handler"
	"ddd-example/internal/user/infra"
)

// Application 账号模块业务逻辑
type Application struct {
	UserRepository adapter.UserRepository

	Authorize         *handler.AuthorizeHandler
	ChangePassword    *handler.ChangePasswordHandler
	LoginWithEmail    *handler.LoginWithEmailHandler
	Logout            *handler.LogoutHandler
	Register          *handler.RegisterHandler
	RegisterWithOauth *handler.RegisterWithOauthHandler
	VerifyOauth       *handler.VerifyOauthHandler
}

// NewApplication 构造函数
func NewApplication(opt *option.Options) *Application {
	db := opt.GetDB()
	cache := infra.NewMemoryCache()

	return initApplication(db, db, cache)
}
