package service

import (
	"context"

	"lmm/api/clock"
	"lmm/api/event"
	"lmm/api/service/user/domain"
	userEvent "lmm/api/service/user/domain/event"
	"lmm/api/service/user/domain/model"
	"lmm/api/testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func TestAssignUserRole(tt *testing.T) {
	type TestCase struct {
		operator    *model.UserDescriptor
		targetUser  *model.UserDescriptor
		targetRole  model.Role
		expectedErr error
	}

	cases := map[string]TestCase{
		"AdminAssginAdminToAdmin": TestCase{
			operator:    newAdmin(),
			targetUser:  newAdmin(),
			targetRole:  model.Admin,
			expectedErr: nil,
		},
		"AdminAssginAdminToOrdinary": TestCase{
			operator:    newAdmin(),
			targetUser:  newAdmin(),
			targetRole:  model.Ordinary,
			expectedErr: nil,
		},
		"AdminAssginOrdinaryToAdmin": TestCase{
			operator:    newAdmin(),
			targetUser:  newOrdinary(),
			targetRole:  model.Admin,
			expectedErr: nil,
		},
		"AdminAssignOrdinaryToOrdinary": TestCase{
			operator:    newAdmin(),
			targetUser:  newOrdinary(),
			targetRole:  model.Ordinary,
			expectedErr: nil,
		},
		"OrdinaryAssignAdminToAdmin": TestCase{
			operator:    newOrdinary(),
			targetUser:  newAdmin(),
			targetRole:  model.Admin,
			expectedErr: domain.ErrNoPermission,
		},
		"OrdinaryAssignAdminToOrdinary": TestCase{
			operator:    newOrdinary(),
			targetUser:  newAdmin(),
			targetRole:  model.Ordinary,
			expectedErr: domain.ErrNoPermission,
		},
		"OrdinaryAssignOrdinaryToAdmin": TestCase{
			operator:    newOrdinary(),
			targetUser:  newOrdinary(),
			targetRole:  model.Admin,
			expectedErr: domain.ErrNoPermission,
		},
		"OrdinaryAssignOrdinaryToOrdinary": TestCase{
			operator:    newOrdinary(),
			targetUser:  newOrdinary(),
			targetRole:  model.Ordinary,
			expectedErr: domain.ErrNoPermission,
		},
		"InvalidPermission": TestCase{
			operator:    newAdmin(),
			targetUser:  newOrdinary(),
			targetRole:  model.Role{},
			expectedErr: domain.ErrNoSuchRole,
		},
	}

	c := context.Background()
	for testname, testcase := range cases {
		tt.Run(testname, func(tt *testing.T) {
			t := testing.NewTester(tt)
			err := AssignUserRole(c, testcase.operator, testcase.targetUser, testcase.targetRole)

			t.Is(testcase.expectedErr, errors.Cause(err))
		})
	}
}

func newAdmin() *model.UserDescriptor {
	return newUserWithRole(model.Admin)
}

func newOrdinary() *model.UserDescriptor {
	return newUserWithRole(model.Ordinary)
}

func newUserWithRole(role model.Role) *model.UserDescriptor {
	randomUserName := func() string {
		return "u" + uuid.New().String()[:7]
	}

	user, err := model.NewUserDescriptor(randomUserName(), role, clock.Now())
	if err != nil {
		panic(err)
	}

	return user
}

func init() {
	event.SyncBus().Subscribe(&userEvent.UserRoleChanged{}, event.NopEventHandler)
}
