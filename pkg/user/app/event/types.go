package event

import "ddd-example/pkg/user/domain"

// Login 账号登录
type Login struct {
	User *domain.User
}

// Register 账号注册
type Register struct {
	User *domain.User
}
