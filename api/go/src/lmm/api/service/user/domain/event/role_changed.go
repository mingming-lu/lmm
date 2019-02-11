package event

import (
	"lmm/api/clock"
	"time"
)

const (
	TopicUserRoleChanged = "user.role.changed"
)

type UserRoleChanged struct {
	operatorUser string
	targetUser   string
	targetRole   string
	changedAt    time.Time
}

func NewUserRoleChangedEvent(operator, user, targetRole string) *UserRoleChanged {
	return &UserRoleChanged{
		operatorUser: operator,
		targetUser:   user,
		targetRole:   targetRole,
		changedAt:    clock.Now(),
	}
}

func (event *UserRoleChanged) Topic() string {
	return TopicUserRoleChanged
}

func (event *UserRoleChanged) OccurredAt() time.Time {
	return event.changedAt
}

func (event *UserRoleChanged) OperatorUser() string {
	return event.operatorUser
}

func (event *UserRoleChanged) TargetUser() string {
	return event.targetUser
}

func (event *UserRoleChanged) TargetRole() string {
	return event.targetRole
}
