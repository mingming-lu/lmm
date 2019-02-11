package messaging

import (
	"context"
	"database/sql"

	"lmm/api/messaging"
	"lmm/api/service/user/domain"
	userEvent "lmm/api/service/user/domain/event"
	"lmm/api/storage/db"
	"lmm/api/util/uuidutil"

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
func (s *Subscriber) OnUserRoleChanged(c context.Context, e messaging.Event) error {
	userRoleChanged, ok := e.(*userEvent.UserRoleChanged)
	if !ok {
		return errors.Wrap(messaging.ErrInvalidEvent, e.Topic())
	}

	tx, err := s.db.Begin(c, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return err
	}

	searchUsers := tx.Prepare(c, `
		select id, role from user where name in (?, ?) order by field (name, ?, ?) for update
	`)
	defer searchUsers.Close()

	setRole := tx.Prepare(c, `
		update user set role = ? where name = ?
	`)
	defer setRole.Close()

	recordChangeHistory := tx.Prepare(c, `
		insert into user_role_change_history (
			operator, operator_role, target_user, from_role, to_role, changed_at) values (?, ?, ?, ?, ?, ?)
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
			return db.RollbackWithError(tx, err)
		}

		rows.Next()
		if err := rows.Scan(&operatorUserID, &operatorUserRole); err != nil {
			return db.RollbackWithError(tx, err)
		}
		rows.Next()
		if err := rows.Scan(&targetUserID, &targetUserRole); err != nil {
			return db.RollbackWithError(tx, err)
		}
		rows.Close()
	}

	{
		_, err := setRole.Exec(c, userRoleChanged.TargetRole(), userRoleChanged.TargetUser())
		if err != nil {
			return db.RollbackWithError(tx, err)
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
			return db.RollbackWithError(tx, err)
		}
	}

	return tx.Commit()
}

// OnUserPasswordChanged implements event handler to handle UserPasswordChanged
func (s *Subscriber) OnUserPasswordChanged(c context.Context, e messaging.Event) error {
	userPasswordChanged, ok := e.(*userEvent.UserPasswordChanged)
	if !ok {
		return errors.Wrap(messaging.ErrInvalidEvent, e.Topic())
	}

	tx, err := s.db.Begin(c, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}

	lockUser := tx.Prepare(c, `select id from user where name = ? for update`)
	defer lockUser.Close()

	recordChanged := tx.Prepare(c,
		`insert into user_password_change_history (user, changed_at) values(?, ?)`,
	)
	defer recordChanged.Close()

	updateUser := tx.Prepare(c, `update user set token = ?, password = ? where id = ?`)
	defer updateUser.Close()

	var userID int64
	if err := lockUser.QueryRow(c, userPasswordChanged.UserName()).Scan(&userID); err != nil {
		if err == sql.ErrNoRows {
			return db.RollbackWithError(tx, domain.ErrNoSuchUser)
		}
		return db.RollbackWithError(tx, err)
	}

	if _, err := recordChanged.Exec(c, userID, userPasswordChanged.OccurredAt()); err != nil {
		return db.RollbackWithError(tx, err)
	}

	newToken := uuidutil.NewUUID()
	password := userPasswordChanged.Password()
	if _, err := updateUser.Exec(c, newToken, password, userID); err != nil {
		return db.RollbackWithError(tx, err)
	}

	return tx.Commit()
}
