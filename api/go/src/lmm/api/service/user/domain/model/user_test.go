package model

import (
	"lmm/api/clock"
	"lmm/api/messaging"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/event"
	"lmm/api/testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func TestAssignUserRole(tt *testing.T) {
	type TestCase struct {
		operator    *User
		targetUser  *User
		targetRole  Role
		expectedErr error
	}

	cases := map[string]TestCase{
		"AdminAssginAdminToAdmin": TestCase{
			operator:    newAdmin(),
			targetUser:  newAdmin(),
			targetRole:  Admin,
			expectedErr: nil,
		},
		"AdminAssginAdminToOrdinary": TestCase{
			operator:    newAdmin(),
			targetUser:  newAdmin(),
			targetRole:  Ordinary,
			expectedErr: nil,
		},
		"AdminAssginOrdinaryToAdmin": TestCase{
			operator:    newAdmin(),
			targetUser:  newOrdinary(),
			targetRole:  Admin,
			expectedErr: nil,
		},
		"AdminAssignOrdinaryToOrdinary": TestCase{
			operator:    newAdmin(),
			targetUser:  newOrdinary(),
			targetRole:  Ordinary,
			expectedErr: nil,
		},
		"OrdinaryAssignAdminToAdmin": TestCase{
			operator:    newOrdinary(),
			targetUser:  newAdmin(),
			targetRole:  Admin,
			expectedErr: domain.ErrNoPermission,
		},
		"OrdinaryAssignAdminToOrdinary": TestCase{
			operator:    newOrdinary(),
			targetUser:  newAdmin(),
			targetRole:  Ordinary,
			expectedErr: domain.ErrNoPermission,
		},
		"OrdinaryAssignOrdinaryToAdmin": TestCase{
			operator:    newOrdinary(),
			targetUser:  newOrdinary(),
			targetRole:  Admin,
			expectedErr: domain.ErrNoPermission,
		},
		"OrdinaryAssignOrdinaryToOrdinary": TestCase{
			operator:    newOrdinary(),
			targetUser:  newOrdinary(),
			targetRole:  Ordinary,
			expectedErr: domain.ErrNoPermission,
		},
		"InvalidPermission": TestCase{
			operator:    newAdmin(),
			targetUser:  newOrdinary(),
			targetRole:  Role{},
			expectedErr: domain.ErrNoSuchRole,
		},
	}

	for testname, testcase := range cases {
		tt.Run(testname, func(tt *testing.T) {
			t := testing.NewTester(tt)

			err := testcase.operator.AssignRole(testcase.targetUser, testcase.targetRole)

			t.Is(testcase.expectedErr, errors.Cause(err))
		})
	}
}

func newAdmin() *User {
	return newUserWithRole(Admin)
}

func newOrdinary() *User {
	return newUserWithRole(Ordinary)
}

func newUserWithRole(role Role) *User {
	randomUserName := "u" + uuid.New().String()[:7]
	email := randomUserName + "@lmm.local"
	password := uuid.New().String()
	token := uuid.New().String()

	user, err := NewUser(randomUserName, email, password, token, role, clock.Now())
	if err != nil {
		panic(err)
	}

	return user
}

func init() {
	messaging.SyncBus().Subscribe(&event.UserRoleChanged{}, messaging.NopEventHandler)
}
