package uploader

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"lmm/api/messaging/rabbitmq"
)

var (
	patternImageName = regexp.MustCompile(`^\w+.(\w+)$`)
)

type RabbitMQAssetUploader struct {
	client *rabbitmq.Client
}

func NewRabbitMQAssetUploader(client *rabbitmq.Client) *RabbitMQAssetUploader {
	return &RabbitMQAssetUploader{
		client: client,
	}
}

func (uploader *RabbitMQAssetUploader) Upload(c context.Context, name string, data []byte, args ...interface{}) error {
	if len(args) != 1 {
		return ErrImageUploadTypeNotGiven
	}

	config, ok := args[0].(ImageUploaderConfig)
	if !ok {
		return ErrImageUploadTypeNotGiven
	}

	switch config.Type {
	case "image", "photo":
	default:
		return errors.Wrap(ErrInvalidImageUploadType, "invalid asset upload type: "+config.Type)
	}

	matched := patternImageName.FindStringSubmatch(name)
	if len(matched) != 2 {
		return errors.Wrap(ErrInvalidAssetName, "invalid name: "+name)
	}

	return uploader.client.Publish(c, rabbitmq.NewEvent(
		&amqp.Publishing{
			MessageId: name,
			Type:      fmt.Sprintf("asset.%s.uploaded", config.Type),
			Timestamp: time.Now(),
			Body:      data,
		},
	))
}

func (uploader *RabbitMQAssetUploader) Delete(c context.Context, name string, args ...interface{}) error {
	return errors.New("not implemented")
}

func (uploader *RabbitMQAssetUploader) Close() error {
	return uploader.client.Close()
}
