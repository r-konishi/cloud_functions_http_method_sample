package user

import (
	"context"
	"log"
)

const (
	// UserCollection is user collection
	UserCollection string = "Users"
)

// User is one user information
type User struct {
	UUID  string `json:"uuid"`
	Count int8   `json:"count"`
}

// Get is a function to get user information
func (u *User) Get(c context.Context) (err error) {
	// 最終的に firestore の client を明示的に Close する
	// と、連続でアクセスした場合にエラーになることがある
	// defer client.Close()

	log.Printf("UserModel Get: %+v", u)
	ds, err := client.Collection(UserCollection).Doc(u.UUID).Get(c)
	if err != nil {
		log.Printf("UserModel Get Error: %+v", err)
		return err
	}
	return ds.DataTo(u)
}

// Create is a function to create user information
func (u *User) Create(c context.Context) (err error) {
	// 最終的に firestore の client を明示的に Close する
	// と、連続でアクセスした場合にエラーになることがある
	// defer client.Close()

	log.Printf("UserModel Create: %+v", u)
	if _, err = client.Collection(UserCollection).Doc(u.UUID).Set(c, u); err != nil {
		// error
		log.Printf("UserModel Create Error: %+v", err)
		return err
	}
	return nil
}
