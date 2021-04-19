package main

import (
    "fmt"
    "database/sql"
    "os"

    _ "github.com/lib/pq"
)

type Table struct {
    Name   string
    Fields string
}

func CreateTables(db *sql.DB) {
    tables := []Table{
        Table{Name: "Accounts",                Fields: "Id serial PRIMARY KEY, Login text, Password text"},
        Table{Name: "Articles",                Fields: "Id text PRIMARY KEY, Title text, Abstract text, LastUpdateTimestamp bigint, FullDocumentURL text"},
        Table{Name: "AuthorsOfArticles",       Fields: "ArticleId text REFERENCES Articles (Id), AuthorName text"},
        Table{Name: "AccountArticleRelations", Fields: "UserId integer REFERENCES Accounts (Id), ArticleId text REFERENCES Articles (Id), IsSubscribed boolean, LastAccess bigint"},
        Table{Name: "AccountSearchRelations",  Fields: "UserId integer REFERENCES Accounts (Id), Search text, IsSubscribed boolean, LastAccess bigint"},
    }

    fmt.Println("Creating tables")

    for _, table := range(tables) {
        _, err := db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", table.Name, table.Fields))
        if err != nil {
            panic(err)
        }
    }
    fmt.Println("Created all tables")
}

func main() {
    fmt.Println("Connecting to DB")
    db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable", os.Getenv("dbusername"), os.Getenv("dbname")))
    if err != nil {
        panic(err)
    }
    defer db.Close()

    CreateTables(db)
}
