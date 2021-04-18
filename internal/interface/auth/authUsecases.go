package auth

import (
    "github.com/mp-hl-2021/unarXiv/internal/domain/model"
    "github.com/mp-hl-2021/unarXiv/internal/interface/accounts"
    "github.com/mp-hl-2021/unarXiv/internal/usecases"

    "golang.org/x/crypto/bcrypt"

    "errors"
    "unicode"
)

type Usecases struct {
    accountRepo accounts.Interface
    auth        Interface
}

func NewUsecases(accountRepo accounts.Interface, auth Interface) *Usecases {
    return &Usecases{
        accountRepo: accountRepo,
        auth:        auth,
    }
}

var (
    ErrInvalidLoginString    = errors.New("login string contains invalid character")
    ErrInvalidPasswordString = errors.New("password string contains invalid character")
    ErrTooShortLogin = errors.New("too short login")
    ErrTooLongLogin  = errors.New("too long login")
    ErrTooShortPassword = errors.New("too short password")
    ErrTooLongPassword = errors.New("too long password")
)

const (
    minLoginLength    = 6
    maxLoginLength    = 20
    minPasswordLength = 14
    maxPasswordLength = 48
)

func validateLogin(login string) error {
    chars := 0
    for _, r := range login {
        if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
            return ErrInvalidLoginString
        }
        chars++
    }
    if chars < minLoginLength {
        return ErrTooShortLogin
    }
    if chars > maxLoginLength {
        return ErrTooLongLogin
    }
    return nil
}

func validatePassword(password string) error {
    chars := 0
    for _, r := range password {
        if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
            return ErrInvalidPasswordString
        }
        chars++
    }
    if chars < minPasswordLength {
        return ErrTooShortPassword
    }
    if chars > maxPasswordLength {
        return ErrTooLongPassword
    }
    return nil
}

func (d *Usecases) Register(request usecases.AuthRequest) (usecases.AuthToken, error) {
    if err := validateLogin(request.Login); err != nil {
        return "", err
    }
    if err := validatePassword(request.Password); err != nil {
        return "", err
    }
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    acc, err := d.accountRepo.CreateAccount(accounts.Credentials{
        Login:    request.Login,
        Password: string(hashedPassword),
    })
    if err != nil {
        return "", err
    }
    token, err := d.auth.IssueToken(acc.Id)
    if err != nil {
        return "", err
    }
    return usecases.AuthToken(token), nil
}

func (d *Usecases) Login(request usecases.AuthRequest) (usecases.AuthToken, error) {
    if err := validateLogin(request.Login); err != nil {
        return "", err
    }
    if err := validatePassword(request.Password); err != nil {
        return "", err
    }
    acc, err := d.accountRepo.GetAccountByLogin(request.Login)
    if err != nil {
        return "", err
    }
    if err := bcrypt.CompareHashAndPassword([]byte(acc.Credentials.Password), []byte(request.Password)); err != nil {
        return "", err
    }
    token, err := d.auth.IssueToken(acc.Id)
    if err != nil {
        return "", err
    }
    return usecases.AuthToken(token), nil
}

func (d *Usecases) Decode(token usecases.AuthToken) (model.UserId, error) {
    id, err := d.auth.UserIdByToken(string(token))
    if err != nil {
        return "", err
    }
    return model.UserId(id), nil
}
