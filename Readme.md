## Rabbitmq-go-consumer

you'll never worry about connection lost when receive message from rabbitmq

### Features

- Reconnect rabbitmq when network error or queue service breakdown
- Use callback goroutine very simple when receiving message
- Recover consumers after auto reconnect

### Installation

```
go get github.com/bryant24/rabbitmq-go-consumer
```




### Quick Start

```
package main

import (
	"fmt"
	rgc "github.com/bryant24/rabbitmq-go-consumer"
)

func main() {
	forever := make(chan bool)

	mqSession := rgc.NewQueue("amqp://admin:admin@192.168.1.100:5672/render_mq", "Test")
	tc := rgc.NewTaskConsumer(mqSession)
	tc.Callback = func(msg string) {
		fmt.Println("This is the go routine job handler", msg)
	}
	go tc.StartConsumer()

	<-forever
}



```
