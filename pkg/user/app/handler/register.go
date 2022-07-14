package handler

import (
	"context"
	"fmt"

	"gitlab.haochang.tv/yangyi/examine-code/pkg/user/app/service"
	"gitlab.haochang.tv/yangyi/examine-code/pkg/user/domain"
)

// Register 账号注册，参数
type Register struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required"`
}

// RegisterHandler 账号注册
type RegisterHandler struct {
	Session *service.SessionTokenService
	Users   *service.UserService
}

// Handle 执行账号注册
func (h *RegisterHandler) Handle(ctx context.Context, args Register) (user *domain.User, token string, err error) {
	user, err = h.Users.Create(ctx, args.Email, args.Password)
	if err != nil {
		return
	}

	token, err = h.Session.Generate(ctx, user)
	if err != nil {
		err = fmt.Errorf("generate session token, %w", err)
	}
	return
}
