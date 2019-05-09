package rabbitmq_go_consumer

import "github.com/streadway/amqp"

type ConsumerFunc func(*amqp.Delivery)

type Consumer struct {
	autoAck  bool
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

func (c *Consumer) AutoAck(auto bool) {
	c.queue.autoAck = auto
}

func (c *Consumer) Ack(ack bool) {

}

func (c *Consumer) StartConsumer() {
	c.queue.Consume(func(d amqp.Delivery) {
		go c.Callback(&d)
	})
}

func (c *Consumer) stopConsumer() {
	c.queue.Close()
}

func (c *Consumer) recoverConsumer() {
	c.queue = NewQueue(c.queue.Url, c.queue.Name)
	c.StartConsumer()
}
