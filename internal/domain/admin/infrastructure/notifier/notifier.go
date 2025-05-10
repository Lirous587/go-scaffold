package notifier

type Notifier interface {
	SendMockNotification(to string, id int) error
}
