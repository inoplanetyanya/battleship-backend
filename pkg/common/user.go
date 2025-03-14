package common

import "fmt"

type User struct {
	Id       int    `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u User) String() string {
	return fmt.Sprintf("User{Id: %d, Username: %s}", u.Id, u.Username)
}

func (u User) Equals(other User) bool {
	return u.Id == other.Id && u.Username == other.Username
}
