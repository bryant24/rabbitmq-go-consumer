package rabbitmq_go_consumer

import (
	"fmt"
	"testing"
)

func TestNewTaskConsumer(t *testing.T) {
	forever := make(chan bool)

	mqSession := NewQueue("amqp://admin:admin@192.168.1.100:5672/render_mq", "Test")
	tc := NewTaskConsumer(mqSession)
	tc.Callback = func(msg string) {
		fmt.Println("This is the go routine job handler", msg)
	}
	go tc.StartConsumer()

	<-forever

}
