package rabbitmq_go_consumer

type ConsumerFunc func(string)

type Consumer struct {
	stopChan chan bool
	closed   bool
	queue    *Mqueue
	Callback ConsumerFunc
}

func NewTaskConsumer(queue *Mqueue) *Consumer {
	consumer := new(Consumer)
	consumer.queue = queue

	return consumer
}

func (c *Consumer) StartConsumer() {
	c.queue.Consume(func(msg string) {
		go c.Callback(msg)
	})
}

func (c *Consumer) stopConsumer() {
	c.queue.Close()
}

func (c *Consumer) recoverConsumer() {
	c.queue = NewQueue(c.queue.Url, c.queue.Name)
	c.StartConsumer()
}
