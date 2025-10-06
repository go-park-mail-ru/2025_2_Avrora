package handlers

type IAuthUsecase interface {
	Register(email string, password string) error
	Login(email string, password string) (string, error)
	Logout() (string, error)
}

type authHandler struct {
	authUsecase IAuthUsecase
}

func NewAuthHandler(uc IAuthUsecase) *authHandler {
	return &authHandler{authUsecase: uc}
}
