package messaging

import (
	"context"
	"database/sql"

	"lmm/api/event"
	"lmm/api/http"
	userEvent "lmm/api/service/user/domain/event"
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
func (s *Subscriber) OnUserRoleChanged(c context.Context, e event.Event) error {
	userRoleChanged, ok := e.(*userEvent.UserRoleChanged)
	if !ok {
		return errors.Wrap(event.ErrInvalidEvent, e.Topic())
	}

	tx, err := s.db.Begin(c, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return err
	}

	searchUsers, err := tx.Prepare(c, `
		select id, role from user where name in (?, ?) order by field (name, ?, ?) for update
	`)
	if err != nil {
		tx.Rollback()
		http.Log().Panic(c, err.Error())
	}
	defer searchUsers.Close()

	setRole, err := tx.Prepare(c, `
		update user set role = ? where name = ?
	`)
	if err != nil {
		tx.Rollback()
		http.Log().Panic(c, err.Error())
	}
	defer setRole.Close()

	recordChangeHistory, err := tx.Prepare(c, `
		insert into user_role_change_history (
			operator, operator_role, target_user, from_role, to_role, changed_at
		) values (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		http.Log().Panic(c, err.Error())
	}
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
			tx.Rollback()
			return err
		}

		rows.Next()
		if err := rows.Scan(&operatorUserID, &operatorUserRole); err != nil {
			tx.Rollback()
			return err
		}
		rows.Next()
		if err := rows.Scan(&targetUserID, &targetUserRole); err != nil {
			tx.Rollback()
			return err
		}
		rows.Close()
	}

	{
		_, err := setRole.Exec(c, userRoleChanged.TargetRole(), userRoleChanged.TargetUser())
		if err != nil {
			tx.Rollback()
			return err
		}
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
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
