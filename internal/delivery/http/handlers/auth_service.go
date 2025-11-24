package handlers

import "context"

type IAuthService interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context) (string, error)
}
