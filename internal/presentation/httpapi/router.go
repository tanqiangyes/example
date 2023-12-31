package httpapi

import (
	"ddd-example/internal/option"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

func newRouter(opt *option.Options) chi.Router {
	router := chi.NewRouter()

	router.Use(recoverer(slog.Default()))

	uc := newUserController(opt)
	router.Use(uc.Authorize)

	router.Post(`/session`, uc.LoginWithEmail())
	router.Post(`/register`, uc.Register())
	router.Delete(`/session`, uc.Logout())
	router.Get(`/login/oauth/{site}`, uc.LoginWithOauth())
	router.Post(`/login/oauth/{site}`, uc.VerifyOauth())
	router.Post(`/register/oauth`, uc.RegisterWithOauth())

	router.Group(func(router chi.Router) {
		router.Use(uc.DenyAnonymous)

		router.Get(`/session`, uc.MyIdentity())
		router.Put(`/my/password`, uc.ChangePassword())
	})

	return router
}
