package memory

import (
    "fmt"
    "github.com/mp-hl-2021/unarXiv/internal/domain"
    "github.com/mp-hl-2021/unarXiv/internal/domain/model"
    "sync"
)

type UserRepo struct {
    idToUser    map[model.UserId]model.User
    loginToUser map[string]model.User
    mutex       *sync.Mutex
}

func NewUserRepo() *UserRepo {
    return &UserRepo{
        idToUser:    make(map[model.UserId]model.User),
        loginToUser: make(map[string]model.User),
        mutex:       &sync.Mutex{},
    }
}

func (u *UserRepo) UserById(id model.UserId) (model.User, error) {
    u.mutex.Lock()
    defer u.mutex.Unlock()
    if user, ok := u.idToUser[id]; !ok {
        return model.User{}, domain.UserNotFound
    } else {
        return user, nil
    }
}

func (u *UserRepo) UserByLogin(login string) (model.User, error) {
    u.mutex.Lock()
    defer u.mutex.Unlock()
    if user, ok := u.loginToUser[login]; !ok {
        return model.User{}, domain.UserNotFound
    } else {
        return user, nil
    }
}

func (u *UserRepo) Register(login string) (model.User, error) {
    u.mutex.Lock()
    defer u.mutex.Unlock()
    if _, ok := u.loginToUser[login]; ok {
        return model.User{}, domain.LoginIsAlreadyTaken
    }
    user := model.User{
        Id:    model.UserId(fmt.Sprint("%v", len(u.loginToUser))),
        Login: login,
    }
    u.loginToUser[login] = user
    u.idToUser[user.Id] = user
    return user, nil
}
