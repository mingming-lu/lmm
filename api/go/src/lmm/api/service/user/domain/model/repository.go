package model

import (
	"lmm/api/pkg/transaction"
)

// UserRepository interface
type UserRepository interface {
	NextID(tx transaction.Transaction) (UserID, error)
	Save(tx transaction.Transaction, user *User) error
	FindByName(tx transaction.Transaction, username string) (*User, error)
	FindByToken(tx transaction.Transaction, token string) (*User, error)
}
