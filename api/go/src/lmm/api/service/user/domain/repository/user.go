package repository

import (
	"lmm/api/pkg/transaction"

	"lmm/api/service/user/domain/model"
)

// UserRepository interface
type UserRepository interface {
	NextID(tx transaction.Transaction) (model.UserID, error)
	Save(tx transaction.Transaction, user *model.User) error
	FindByName(tx transaction.Transaction, username string) (*model.User, error)
	// DescribeAll(tx transaction.Transaction, options DescribeAllOptions) ([]*model.UserDescriptor, uint, error)
}

type DescribeAllOptions struct {
	Page  uint
	Count uint
	Order DescribeAllOrder
}

type DescribeAllOrder int

const (
	DescribeAllOrderByNameAsc = iota
	DescribeAllOrderByNameDesc
	DescribeAllOrderByRegisteredDateAsc
	DescribeAllOrderByRegisteredDateDesc
	DescribeAllOrderByRoleAsc
	DescribeAllOrderByRoleDesc
)
