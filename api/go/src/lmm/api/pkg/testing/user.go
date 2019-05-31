package testing

import "cloud.google.com/go/datastore"

// User is expected to be useful for testing
type User struct {
	id       int64
	name     string
	email    string
	password string
}

// NewUser randomly create a new user
func NewUser(source *datastore.Client) *User {
	return &User{}
}
