package messaging

import (
	"context"
	"testing"

	"lmm/api/messaging"
	"lmm/api/pkg/pubsub"
	"lmm/api/pkg/pubsub/pubsubtest"
	"lmm/api/service/user/domain/model"

	"github.com/stretchr/testify/assert"
)

func TestUserEventPublisher(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		UserID     model.UserID
		AckMsg     string
		NotifyFunc func(pub model.UserEventPublisher) func(context.Context, model.UserID) error
	}{
		TopicUserPasswordChanged: {
			UserID: model.UserID(123),
			AckMsg: "password changed",
			NotifyFunc: func(pub model.UserEventPublisher) func(context.Context, model.UserID) error {
				return pub.NotifyUserPasswordChanged
			},
		},
		TopicUserRegistered: {
			UserID: model.UserID(541),
			AckMsg: "registered",
			NotifyFunc: func(pub model.UserEventPublisher) func(context.Context, model.UserID) error {
				return pub.NotifyUserRegistered
			},
		},
	}

	for topic, testcase := range cases {
		t.Run(topic, func(t *testing.T) {
			sigChan := make(chan string, 1)

			client := pubsubtest.NewClient()
			go client.Subscribe(ctx, topic, func(c context.Context, evt messaging.Event) error {
				var actual userEvent

				assert.Equal(t, topic, evt.Topic())
				assert.NoError(t, pubsub.ScanEvent(evt, &actual))
				assert.Equal(t, int(testcase.UserID), actual.UserID)

				sigChan <- testcase.AckMsg
				return nil
			})

			pub := NewUserEventPublisher(client)
			assert.NoError(t, testcase.NotifyFunc(pub)(ctx, testcase.UserID))
			assert.Equal(t, testcase.AckMsg, <-sigChan)

			client.Close()
		})
	}

}
