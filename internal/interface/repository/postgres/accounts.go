package postgres

import (
	"database/sql"
	"github.com/mp-hl-2021/unarXiv/internal/interface/accounts"
	"strconv"

	_ "github.com/lib/pq"
)

type AccountsRepo struct {
	db *sql.DB
}

func NewAccountsRepo(db *sql.DB) *AccountsRepo {
	return &AccountsRepo{db: db}
}

func (m *AccountsRepo) GetAccountById(id string) (accounts.Account, error) {
	rows, err := m.db.Query("SELECT * FROM Accounts where id=$1;", id)
	if err != nil {
		return accounts.Account{}, err
	}
	defer rows.Close()
	a := accounts.Account{}
	for rows.Next() {
		err := rows.Scan(&a.Id, &a.Credentials.Login, &a.Credentials.Password)
		return a, err
	}
	return a, accounts.ErrNotFound
}

func (m *AccountsRepo) GetAccountByLogin(login string) (accounts.Account, error) {
	rows, err := m.db.Query("SELECT * FROM Accounts where login=$1;", login)
	if err != nil {
		return accounts.Account{}, err
	}
	defer rows.Close()
	a := accounts.Account{}
	for rows.Next() {
		err := rows.Scan(&a.Id, &a.Credentials.Login, &a.Credentials.Password)
		return a, err
	}
	return a, accounts.ErrNotFound
}

func (m *AccountsRepo) CreateAccount(cred accounts.Credentials) (accounts.Account, error) {
	if _, err := m.GetAccountByLogin(cred.Login); err == nil {
		return accounts.Account{}, accounts.ErrAlreadyExists
	}
	var id uint64
	err := m.db.QueryRow("INSERT INTO Accounts (Login, Password) VALUES ($1, $2) RETURNING Id;", cred.Login, cred.Password).Scan(&id)
	if err != nil {
		return accounts.Account{}, err
	}
	return accounts.Account{
		Id:          strconv.FormatUint(id, 16),
		Credentials: cred,
	}, nil
}
