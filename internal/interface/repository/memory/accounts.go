package memory

import (
    "github.com/mp-hl-2021/unarXiv/internal/interface/accounts"
    "strconv"
    "fmt"
    "database/sql"

    _ "github.com/lib/pq"
)

type AccountsRepo struct {
    db *sql.DB
}

func NewAccountsRepo(db *sql.DB) *AccountsRepo {
    return &AccountsRepo{db: db}
}

func (m *AccountsRepo) GetAccountWithCondition(condition string) (accounts.Account, error) {
    rows, err := m.db.Query("SELECT * FROM Accounts where $1;", condition)
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    a := accounts.Account{}
    for rows.Next() {
        if err := rows.Scan(&a.Id, &a.Credentials.Login, &a.Credentials.Password); err != nil {
            panic(err)
        } else {
            return a, nil
        }
    }
    return a, accounts.ErrNotFound
}

func (m *AccountsRepo) GetAccountById(id string) (accounts.Account, error) {
    return m.GetAccountWithCondition(fmt.Sprintf("id=%d", id))
}

func (m *AccountsRepo) GetAccountByLogin(login string) (accounts.Account, error) {
    return m.GetAccountWithCondition(fmt.Sprintf("login='%s'", login))
}

func (m *AccountsRepo) CreateAccount(cred accounts.Credentials) (accounts.Account, error) {
    if _, err := m.GetAccountByLogin(cred.Login); err == nil {
        return accounts.Account{}, accounts.ErrAlreadyExists
    }
    var id uint64
    err := m.db.QueryRow("INSERT INTO Accounts (Login, Password) VALUES ($1, $2) RETURNING Id;", cred.Login, cred.Password).Scan(&id)
    if err != nil {
        panic(err)
    }
    return accounts.Account{
        Id:          strconv.FormatUint(id, 16),
        Credentials: cred,
    }, nil
}
