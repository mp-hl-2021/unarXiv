package smartUsecases

import (
    "github.com/mp-hl-2021/unarXiv/internal/domain/model"
    "github.com/mp-hl-2021/unarXiv/internal/usecases"
    "github.com/mp-hl-2021/unarXiv/internal/accountstorage"
    "github.com/mp-hl-2021/unarXiv/internal/interface/auth"

    "golang.org/x/crypto/bcrypt"

	"errors"
	"unicode"
)

var smartToken usecases.AuthToken = "jwt"
var smartArticle = model.ArticleMeta{
    Id:                  "smart",
    Title:               "bunny",
    Authors:             []string{"ya"},
    Abstract:            "abstract",
    LastUpdateTimestamp: 0,
}
var smartArticleSubscription = model.UserArticleSubscription{
    UserId:    0,
    ArticleId: "smart",
}
var smartSearchSubscription = model.UserSearchSubscription{
    UserId: 0,
    Query:  "smart",
}

type SmartUsecases struct {
    AccountStorage accountstorage.Interface
	Auth           auth.Interface
}

var (
	ErrInvalidLoginString    = errors.New("login string contains invalid character")
	ErrInvalidPasswordString = errors.New("password string contains invalid character")
	ErrTooShortString        = errors.New("too short string")
	ErrTooLongString         = errors.New("too long string")

	ErrInvalidLogin    = errors.New("login not found")
	ErrInvalidPassword = errors.New("invalid password")
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
		return ErrTooShortString
	}
	if chars > maxLoginLength {
		return ErrTooLongString
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
		return ErrTooShortString
	}
	if chars > maxPasswordLength {
		return ErrTooLongString
	}
	return nil
}

func (d *SmartUsecases) Register(request usecases.AuthRequest) (usecases.AuthToken, error) {
    if err := validateLogin(request.Login); err != nil {
		return usecases.AuthToken{}, err
	}
	if err := validatePassword(request.Password); err != nil {
		return usecases.AuthToken{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return usecases.AuthToken{}, err
	}
	acc, err := d.AccountStorage.CreateAccount(accountstorage.Credentials{
		Login:    request.Login,
		Password: string(hashedPassword),
	})
	if err != nil {
		return usecases.AuthToken{}, err
	}
	return usecases.AuthToken{acc.Id}, nil
}

func (d *SmartUsecases) Login(request usecases.AuthRequest) (usecases.AuthToken, error) {
    if err := validateLogin(request.Login); err != nil {
		return usecases.AuthToken{}, err
	}
	if err := validatePassword(request.Password); err != nil {
		return usecases.AuthToken{}, err
	}
	acc, err := d.AccountStorage.GetAccountByLogin(request.Login)
	if err != nil {
		return usecases.AuthToken{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(acc.Credentials.Password), []byte(password)); err != nil {
		return usecases.AuthToken{}, err
	}
	token, err := d.Auth.IssueToken(acc.Id)
	if err != nil {
		return usecases.AuthToken{}, err
	}
	return usecases.AuthToken{token}, nil
}

func (d *SmartUsecases) Decode(token usecases.AuthToken) (model.UserId, error) {
    return d.Auth.UserIdByToken()
}

func (d *SmartUsecases) AccessArticle(articleId model.ArticleId, userId *model.UserId) (model.Article, error) {
    return model.Article{
        ArticleMeta: smartArticle,
    }, nil
}

func (d *SmartUsecases) Search(query model.SearchQuery, userId *model.UserId) (model.SearchResult, error) {
    return model.SearchResult{
        TotalMatchesCount: 3,
        Articles: []model.ArticleMeta{smartArticle},
    }, nil
}

func (d *SmartUsecases) GetSearchHistory(id model.UserId) (model.UserSearchHistory, error) {
    return model.UserSearchHistory{
        UserId:  0,
        Queries: []string{"smart"},
    }, nil
}

func (d *SmartUsecases) ClearSearchHistory(id model.UserId) error {
    return nil
}

func (d *SmartUsecases) GetArticleHistory(id model.UserId) (model.UserArticleHistory, error) {
    return model.UserArticleHistory{
        UserId:   0,
        Articles: []model.ArticleMeta{smartArticle},
    }, nil
}

func (d *SmartUsecases) ClearArticleHistory(id model.UserId) error {
    return nil
}

func (d *SmartUsecases) GetArticleLastAccess(userId model.UserId, articleId model.ArticleId) (*model.UserArticleAccess, error) {
    return &model.UserArticleAccess{
        UserId:    0,
        ArticleId: "smart",
        Timestamp: 293,
    }, nil
}

func (d *SmartUsecases) GetSearchLastAccess(userId model.UserId, query string) (*model.UserSearchAccess, error) {
    return &model.UserSearchAccess{
        UserId:    0,
        Query:     "smart",
        Timestamp: 2394,
    }, nil
}

func (d *SmartUsecases) SubscribeForArticle(userId model.UserId, articleId model.ArticleId) (model.UserArticleSubscription, error) {
    return smartArticleSubscription, nil
}

func (d *SmartUsecases) UnsubscribeFromArticle(userId model.UserId, articleId model.ArticleId) error {
    return nil
}

func (d *SmartUsecases) CheckArticleSubscription(userId model.UserId, articleId model.ArticleId) (*model.UserArticleSubscription, error) {
    return &smartArticleSubscription, nil
}

func (d *SmartUsecases) GetArticleSubscriptions(userId model.UserId) ([]model.UserArticleSubscription, error) {
    return []model.UserArticleSubscription{smartArticleSubscription}, nil
}

func (d *SmartUsecases) GetArticleUpdates(userId model.UserId) ([]model.ArticleMeta, error) {
    return []model.ArticleMeta{smartArticle}, nil
}

func (d *SmartUsecases) SubscribeForSearch(userId model.UserId, query string) (model.UserSearchSubscription, error) {
    return smartSearchSubscription, nil
}

func (d *SmartUsecases) UnsubscribeFromSearch(userId model.UserId, query string) error {
    return nil
}

func (d *SmartUsecases) CheckSearchSubscription(userId model.UserId, query string) (*model.UserSearchSubscription, error) {
    return &smartSearchSubscription, nil
}

func (d *SmartUsecases) GetSearchSubscriptions(userId model.UserId) ([]model.UserSearchSubscription, error) {
    return []model.UserSearchSubscription{smartSearchSubscription}, nil
}

func (d *SmartUsecases) GetSearchUpdates(userId model.UserId) ([]string, error) {
    return []string{"smart"}, nil
}

