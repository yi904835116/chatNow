package users

import (
	"errors"
)

//ErrUserNotFound is returned when the user can't be found
var ErrUserNotFound = errors.New("user not found")

//Store represents a store for Users
type Store interface {
	//GetByID returns the User with the given ID
	GetByID(id UserID) (*User, error)

	//GetByEmail returns the User with the given email
	GetByEmail(email string) (*User, error)

	//GetByUserName returns the User with the given Username
	GetByUserName(username string) (*User, error)

	//Insert inserts the user into the database, and returns
	//the newly-inserted User, complete with the newly-assigned UserID
	Insert(newUser *User) (*User, error)

	//Update applies UserUpdates to the given user ID
	//and returns the newly-updated user
	Update(userID UserID, updates *Updates) (*User, error)

	//Delete deletes the user with the given ID
	Delete(userID UserID) error
}
