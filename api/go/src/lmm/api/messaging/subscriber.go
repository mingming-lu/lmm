package messaging

type Subscriber interface {
	Subscribe(e Event, handler EventHandler) error
}
