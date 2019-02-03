package repository

import (
	"context"

	"lmm/api/service/user/domain/model"
)

// UserRepository interface
type UserRepository interface {
	Save(c context.Context, user *model.User) error
	FindByName(c context.Context, username string) (*model.User, error)
	DescribeByName(c context.Context, username string) (*model.UserDescriptor, error)
	DescribeAll(c context.Context, options DescribeAllOptions) ([]*model.UserDescriptor, error)
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
