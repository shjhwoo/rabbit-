package main

import (
	"bytes"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go" //consumer
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

/*
we will use manual message acknowledgements by passing a false for the "auto-ack" argument 
and then send a proper acknowledgment from the worker with d.Ack(false) 
(this acknowledges a single delivery), once we're done with a task.
즉 메세지 전송중에 컨슈머가 죽었다가 이후 다시 살아나도 메세지를 받을 수 있음
*/
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare( 
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto-ack  *****************이 부분이 false로 변경!
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	  )
	  failOnError(err, "Failed to register a consumer")
	  
	  var forever chan struct{}
	  
	  go func() {
		for d := range msgs {
		  log.Printf("Received a message: %s", d.Body)
		  dotCount := bytes.Count(d.Body, []byte("."))
		  t := time.Duration(dotCount)
		  time.Sleep(t * time.Second)
		  log.Printf("Done")
		  d.Ack(false) // **************************이 부분 추가
		}
	  }()
	  
	  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	  <-forever
}

//프로세싱 중인 메세지를 잃어버리지 않기 위해서 컨슈머 측에서는 ack 옵션을 사용할 수 있다. 
