package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string) string {
    var s string
    if (len(args) < 2) || os.Args[1] == "" {
        s = "hello"
    } else {
        s = strings.Join(args[1:], " ")
    }
    return s
}

func main() {
	//rabbitmq server에 연결. 소켓 연결을 abstraction
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// channel 연다. channel은 rabbitmq api가 위치하고 있는 곳이다
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//큐 선언
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable :: rabbitmq 서버 종료시, false값으로 설정되어 있으면 모든 큐(메세지)가 날아가버림.
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//보낼 메세지 생성하여 큐에 저장해 둔다. publish a message to the queue (이후에 컨슈머가 소비하게된다.)
	body := bodyFrom(os.Args)
	err = ch.PublishWithContext(ctx,
	"",     // exchange: 별도의 라우팅 규칙이 없음
	q.Name, // routing key
	false,  // mandatory
	false,  // immediate
	amqp.Publishing {
		DeliveryMode: amqp.Persistent,
		ContentType: "text/plain",
		Body:        []byte(body), //바이트 배열이라 원하는 대로 인코딩할 수 있다.
	})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body) //
}