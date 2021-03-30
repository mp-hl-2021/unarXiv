package usecases

import "github.com/mp-hl-2021/unarXiv/internal/domain/model"

type AuthRequest struct {
    Login    string
    Password string
}

type AuthToken string

type AuthInterface interface {
    Register(request AuthRequest) (AuthToken, error)
    Login(request AuthRequest) (AuthToken, error)

    Decode(token AuthToken) (model.UserId, error)
}
