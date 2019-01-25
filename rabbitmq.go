package rabbitmq_go_consumer

import (
	"github.com/streadway/amqp"
	"log"
)

type messageConsumer func(string)

type Mqueue struct {
	Url          string
	Name         string
	errorChannel chan *amqp.Error
	connection   *amqp.Connection
	channel      *amqp.Channel
	closed       bool

	consumers []messageConsumer
}

func NewQueue(url string, qName string) *Mqueue {
	q := new(Mqueue)
	q.Url = url
	q.Name = qName
	q.consumers = make([]messageConsumer, 0)

	q.connect()
	go q.reconnect()

	return q
}

func (q *Mqueue) Send(message string) {
	err := q.channel.Publish("", // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		logError("Queue declaration failed", err)
	}
}

func (q *Mqueue) Consume(consumer messageConsumer) {
	deliveries, err := q.registerQueueConsumer()
	q.executeMessageConsumer(err, consumer, deliveries, false)
}

func (q *Mqueue) connect() {
	for {
		conn, err := amqp.Dial(q.Url)
		if err == nil {
			q.connection = conn
			q.errorChannel = make(chan *amqp.Error)
			q.connection.NotifyClose(q.errorChannel)

			q.openChannel()
			q.declareQueue()

			return
		}
	}
}

func (q *Mqueue) reconnect() {
	for {
		err := <-q.errorChannel
		if !q.closed {
			logError("Reconnecting after connection closed", err)

			q.connect()
			q.recoverConsumers()
		} else {
			return
		}
	}
}

func (q *Mqueue) declareQueue() {
	_, err := q.channel.QueueDeclare(
		q.Name,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logError("Queue declaration failed", err)
	}

}

func (q *Mqueue) registerQueueConsumer() (<-chan amqp.Delivery, error) {
	msgs, err := q.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logError("Consuming messages from Mqueue failed", err)
	}

	return msgs, err
}

func (q *Mqueue) executeMessageConsumer(err error, consumer messageConsumer, deliveries <-chan amqp.Delivery, isRecovery bool) {
	if err == nil {
		if !isRecovery {
			q.consumers = append(q.consumers, consumer)
		}
		go func() {
			for delivery := range deliveries {
				consumer(string(delivery.Body[:]))
			}
		}()
	}
}

func (q *Mqueue) recoverConsumers() {
	for i := range q.consumers {
		var consumer = q.consumers[i]

		log.Println("Recovering consumer...")
		msgs, err := q.registerQueueConsumer()
		log.Println("Consumer recovered! Continuing message processing...")
		q.executeMessageConsumer(err, consumer, msgs, true)
	}
}

func (q *Mqueue) openChannel() {
	channel, err := q.connection.Channel()
	if err != nil {
		logError("Opening channel failed", err)
		return
	}
	q.channel = channel
}

func (q *Mqueue) Close() {
	q.closed = true
	q.channel.Close()
	q.connection.Close()
}

func logError(message string, err error) {
	log.Printf("%s:%s", message, err)

}
