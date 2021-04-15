package domain

import "fmt"

var (
	UserNotFound = fmt.Errorf("user not found")
	LoginIsAlreadyTaken = fmt.Errorf("login is already taken")

	AlreadySubscribed = fmt.Errorf("already subscribed")
	NotSubscribed = fmt.Errorf("not subscribed")

	NeverAccessed = fmt.Errorf("never accessed")

	ArticleNotFound = fmt.Errorf("article not found")
)
