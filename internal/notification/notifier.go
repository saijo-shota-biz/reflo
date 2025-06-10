package notification

//go:generate mockgen -source=notifier.go -destination=../../mock/notification/notifier_mock.go -package=mock_notification

type Notifier interface {
	NotifyFocusComplete() error
	NotifyBreakComplete() error
}
