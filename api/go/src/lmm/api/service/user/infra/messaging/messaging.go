package messaging

import (
	"context"

	domain "lmm/api/domain/event"
	"lmm/api/service/user/domain/event"
	"lmm/api/storage/db"

	"github.com/pkg/errors"
)

// Subscriber handles UserRoleChanged
type Subscriber struct {
	db db.DB
}

// NewSubscriber creator
func NewSubscriber(db db.DB) *Subscriber {
	return &Subscriber{db}
}

// OnUserRoleChanged implements event handler to handle UserRoleChanged
func (s *Subscriber) OnUserRoleChanged(c context.Context, e domain.Event) error {
	userRoleChanged, ok := e.(*event.UserRoleChanged)
	if !ok {
		return errors.Wrap(domain.ErrInvalidEvent, e.Topic())
	}

	searchUsers := s.db.Prepare(c, `
		select id, role from user where name in (?, ?) order by filed (name, ?, ?)
	`)
	defer searchUsers.Close()

	recordChangeHistory := s.db.Prepare(c, `
		insert into user_role_change_history
		operator, operator_role, target_user, from, to, changed_at
		values (?, ?, ?, ?, ?, ?)
	`)
	defer recordChangeHistory.Close()

	var (
		operatorUserID   int64
		operatorUserRole string
		targetUserID     int64
		targetUserRole   string
	)

	{
		rows, err := searchUsers.Query(c,
			userRoleChanged.OperatorUser(), userRoleChanged.TargetUser(),
			userRoleChanged.OperatorUser(), userRoleChanged.TargetUser(),
		)
		if err != nil {
			return err
		}

		rows.Next()
		if err := rows.Scan(&operatorUserID, &operatorUserRole); err != nil {
			return err
		}
		rows.Next()
		if err := rows.Scan(&targetUserID, &targetUserRole); err != nil {
			return err
		}
		rows.Close()
	}

	{
		_, err := recordChangeHistory.Exec(c,
			operatorUserID,
			operatorUserRole,
			targetUserID,
			targetUserRole,
			userRoleChanged.TargetRole(),
			userRoleChanged.OccurredAt(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
