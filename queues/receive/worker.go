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

//하나의 Queue에 여러 Consumer가 존재할 경우, Queue는 기본적으로 Round-Robin 방식으로 메세지를 분배합니다.
//이를 예방하기 위해, prefetch count를 1로 설정해 두면, 하나의 메세지가 처리되기 전(Ack를 보내기 전)에는 새로운 메세지를 받지 않게 되므로, 작업을 분산시킬 수 있습니다.
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare( //send.go에서 선언해주었던 큐와 일치시켜야 한다 + 또한 메세지 소비 전에 큐를 미리 선언해줘야 메세지 소비가 가능하기 때문이다.
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count :: 이 부분은 1로 설정해 두어야 여러개의 컨슈머가 모두 골고루 일을 할 수 있게된다. 
		//디폴트인 라운드로빈 방식으로만 하면 하나의 컨슈머에만 너무 일이 몰릴 수 있기 때문에, 놀고있는 애한테도 골고루 일을 줘야한다
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
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
		}
	  }()
	  
	  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	  <-forever
}

//프로세싱 중인 메세지를 잃어버리지 않기 위해서 컨슈머 측에서는 ack 옵션을 사용할 수 있다. 
/*
이는 메세지를 받았고, 처리가 완료되었으므로, mq서버가 해당 컨슈머를 삭제해도 된다는 것을 의미한다
=> 다시말하면 여러개의 컨슈머 중 하나가 죽었을 때 mq서버는 메세지를 다른 컨슈머에게 전달하도록 설정하는 옵션이다. 
=> 현재 코드는 위 상황에서 전달한 메세지를 잃어버릴 위험이 있다. 
=> ack 옵션을 설정한 상황을 가정해보자.
만약 컨슈머가 죽어서 ack신호를 못보내면 mq서버는 메세지가 완전히 처리되지 못했다고 판단하고 그것을 다시 큐에 전달해 줄 것이다. 
만약에 다른 컨슈머들도 있다면, mq는 그 컨슈머들한테 메세지를 즉시 보냏 것이다. 이를 통해서 어떠한 메세지도 유실되지 않도록 한다. 
이와 관련해 ack응답이 돌아오는 타임아웃 값을 설정하여, 버그가 있는 컨슈머들을 잡아낼 수 있다
*/