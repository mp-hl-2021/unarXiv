package repository

import "github.com/mp-hl-2021/unarXiv/internal/domain/model"

type UserRepo interface {
    UserById(id model.UserId) (model.User, error)
    UserByLogin(login string) (model.User, error)

    Register(login string) (model.User, error)
}
