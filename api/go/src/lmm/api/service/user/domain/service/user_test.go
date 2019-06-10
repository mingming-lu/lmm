package service

import (
	"context"
	"testing"

	"lmm/api/clock"
	"lmm/api/messaging"
	"lmm/api/service/user/domain"
	userEvent "lmm/api/service/user/domain/event"
	"lmm/api/service/user/domain/model"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAssignUserRole(t *testing.T) {
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

	c := context.Background()
	for testname, testcase := range cases {
		t.Run(testname, func(t *testing.T) {
			err := AssignUserRole(c, testcase.operator, testcase.targetUser, testcase.targetRole)

			assert.Equal(t, testcase.expectedErr, errors.Cause(err))
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
	randomUserName := "u" + uuid.New().String()[:7]
	email := randomUserName + "@lmm.local"
	password := uuid.New().String()
	token := uuid.New().String()

	user, err := model.NewUser(0, randomUserName, email, password, token, role, clock.Now())
	if err != nil {
		panic(err)
	}

	return user
}

func init() {
	messaging.SyncBus().Subscribe(&userEvent.UserRoleChanged{}, messaging.NopEventHandler)
}
