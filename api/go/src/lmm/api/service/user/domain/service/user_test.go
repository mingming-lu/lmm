package service

import (
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func TestAssignUserRole(tt *testing.T) {
	type TestCase struct {
		operator    *model.User
		targetUser  *model.User
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

	for testname, testcase := range cases {
		tt.Run(testname, func(tt *testing.T) {
			t := testing.NewTester(tt)
			err := AssignUserRole(testcase.operator, testcase.targetUser, testcase.targetRole)

			t.Is(testcase.expectedErr, errors.Cause(err))
		})
	}
}

func newAdmin() *model.User {
	return newUserWithRole(model.Admin)
}

func newOrdinary() *model.User {
	return newUserWithRole(model.Ordinary)
}

func newUserWithRole(role model.Role) *model.User {
	randomString := uuid.New().String

	pw, err := model.NewPassword(randomString())
	if err != nil {
		panic(err)
	}

	user, err := model.NewUser(role.Name(), *pw, randomString(), role)
	if err != nil {
		panic(err)
	}

	return user
}
