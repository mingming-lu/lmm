package messaging

import (
	"context"
	"time"

	"lmm/api/messaging"
	"lmm/api/service/user/domain/model"
)

const (
	TopicUserPasswordChanged = "UserPasswordChanged"
	TopicUserRegistered      = "UserRegistered"
)

type userEventPublisher struct {
	client messaging.Publisher
}

func NewUserEventPublisher(pub messaging.Publisher) model.UserEventPublisher {
	return &userEventPublisher{
		client: pub,
	}
}

type userEvent struct {
	UserID int `json:"user_id"`

	topic       string
	publishedAt time.Time
}

func (e *userEvent) Topic() string {
	return e.topic
}

func (e *userEvent) PublishedAt() time.Time {
	return e.publishedAt
}

func (e *userEvent) Message() interface{} {
	return e
}

func (p *userEventPublisher) NotifyUserPasswordChanged(c context.Context, userID model.UserID) error {
	return p.client.Publish(c, &userEvent{
		UserID:      int(userID),
		topic:       TopicUserPasswordChanged,
		publishedAt: time.Now(),
	})
}

func (p *userEventPublisher) NotifyUserRegistered(c context.Context, userID model.UserID) error {
	return p.client.Publish(c, &userEvent{
		UserID:      int(userID),
		topic:       TopicUserRegistered,
		publishedAt: time.Now(),
	})
}
