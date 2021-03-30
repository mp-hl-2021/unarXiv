package accountstorage

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrAlreadyExist = errors.New("already exist")
)

type Account struct {
	Id uint64
	Credentials
}

type Credentials struct {
	Login    string
	Password string
}

type Interface interface {
	CreateAccount(cred Credentials) (Account, error)
	GetAccountById(id uint64) (Account, error)
	GetAccountByLogin(login string) (Account, error)
}
