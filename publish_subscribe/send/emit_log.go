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

func main() {
        conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
        failOnError(err, "Failed to connect to RabbitMQ")
        defer conn.Close()

        ch, err := conn.Channel()
        failOnError(err, "Failed to open a channel")
        defer ch.Close()

        err = ch.ExchangeDeclare(
                "logs",   // name ****************ch.PublishWithContext에 들어갈 이름임
                "fanout", // type *********************** 중요
                true,     // durable
                false,    // auto-deleted
                false,    // internal
                false,    // no-wait
                nil,      // arguments
        )
        failOnError(err, "Failed to declare an exchange")

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        body := bodyFrom(os.Args)
        err = ch.PublishWithContext(ctx,
                "logs", // exchange **********************************
                "",     // routing key
                false,  // mandatory
                false,  // immediate
                amqp.Publishing{
                        ContentType: "text/plain",
                        Body:        []byte(body),
                })
        failOnError(err, "Failed to publish a message")

        log.Printf(" [x] Sent %s", body)
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


// 	//큐 선언
// 	/*
// 	아래와 같이 name을 빈 문자열, exclusive를 true로 설정하면 현재 주고받는 메세지와 모든 메세지를 수신할 수 있다.
// 	*/
// 	q, err := ch.QueueDeclare(
// 		"", // name
// 		true,   // durable :: rabbitmq 서버 종료시, false값으로 설정되어 있으면 모든 큐(메세지)가 날아가버림.
// 		false,   // delete when unused
// 		true,   // exclusive (name을 빈 문자열로 설정하면 랜덤한 큐 생성되는데 이때 이 옵션을 true로 해줘야 한다.연결 삭제시 큐도 같이 삭제하는 옵션임)
// 		false,   // no-wait
// 		nil,     // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")


//위 파일에서는 별도로 큐를 선언하지 않았지만 랜덤한 이름의 큐가 자동으로 생성된다.

