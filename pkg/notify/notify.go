package notify

type Notifier interface {
	SendSuccessEvent(title, message string) error
	SendFailEvent(title, message string) error
}

func NewNotifier(endpoint, configPath string) (Notifier, error) {
	switch endpoint {
	case "slack":
		sn, err := NewSlackNotify(configPath)
		return &sn, err
	default:
		return &DummyNotifier{}, nil
	}
}
