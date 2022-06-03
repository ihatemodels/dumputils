package notifiers

type Notifier interface {
	Init() error
	Send() error
}
