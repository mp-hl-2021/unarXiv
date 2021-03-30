package auth

import (
    "github.com/mp-hl-2021/unarXiv/internal/domain/model"
    "github.com/mp-hl-2021/unarXiv/internal/usecases"
)

type Interface interface {
	IssueToken(userId model.UserId) (usecases.AuthToken, error)
	UserIdByToken(token usecases.AuthToken) (model.UserId, error)
}
