package memory

import (
	"github.com/mp-hl-2021/unarXiv/internal/interface/accounts"
	"strconv"
	"sync"
)

type AccountsRepo struct {
	accountsById    map[string]accounts.Account
	accountsByLogin map[string]accounts.Account
	nextId          uint64
	mutex           *sync.Mutex
}

func NewAccountsRepo() *AccountsRepo {
	return &AccountsRepo{
		accountsById:    make(map[string]accounts.Account),
		accountsByLogin: make(map[string]accounts.Account),
		mutex:           &sync.Mutex{},
	}
}

func (m *AccountsRepo) CreateAccount(cred accounts.Credentials) (accounts.Account, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.accountsByLogin[cred.Login]; ok {
		return accounts.Account{}, accounts.ErrAlreadyExists
	}
	a := accounts.Account{
		Id:          strconv.FormatUint(m.nextId, 16),
		Credentials: cred,
	}
	m.accountsById[a.Id] = a
	m.accountsByLogin[a.Login] = a
	m.nextId++
	return a, nil
}

func (m *AccountsRepo) GetAccountById(id string) (accounts.Account, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	a, ok := m.accountsById[id]
	if !ok {
		return a, accounts.ErrNotFound
	}
	return a, nil
}

func (m *AccountsRepo) GetAccountByLogin(login string) (accounts.Account, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	a, ok := m.accountsByLogin[login]
	if !ok {
		return a, accounts.ErrNotFound
	}
	return a, nil
}
