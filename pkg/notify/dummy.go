package notify

import (
	"fmt"
)

type DummyNotifier struct{}

func (d *DummyNotifier) SendSuccessEvent(title, message string) error {
	fmt.Printf("title: %s\n", title)
	fmt.Printf("message: %s\n", message)
	fmt.Println("status: Success Event")
	return nil
}

func (d *DummyNotifier) SendFailEvent(title, message string) error {
	fmt.Printf("title: %s\n", title)
	fmt.Printf("message: %s\n", message)
	fmt.Println("status: Fail Event")
	return nil
}
