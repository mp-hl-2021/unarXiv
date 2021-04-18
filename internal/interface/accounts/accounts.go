package accounts

import "errors"

var (
    ErrNotFound      = errors.New("not found")
    ErrAlreadyExists = errors.New("already exists")
)

type Account struct {
    Id string
    Credentials
}

type Credentials struct {
    Login    string
    Password string
}

type Interface interface {
    CreateAccount(cred Credentials) (Account, error)
    GetAccountById(id string) (Account, error)
    GetAccountByLogin(login string) (Account, error)
}
