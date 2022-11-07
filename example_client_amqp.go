package gocelery

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"
)

func ExampleClientAMQP() {
	host_url := "amqp://localhost:5672/"
	amqp_broker := NewAMQPCeleryBroker(host_url)
	amqp_backend := NewAMQPCeleryBackend(host_url)
	cli, err := NewCeleryClient(amqp_broker, amqp_backend, 1)
	if err != nil {
		panic(err)
	}
	cli.SetVersion(true)

	taskName := "worker.add"
	argA := rand.Intn(10)
	argB := rand.Intn(10)

	fmt.Println("Sending message.")
	ayncResult, err := cli.DelayKwargs(
		taskName,
		map[string]interface{}{
			"a": argA,
			"b": argB,
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Message sent.")
	fmt.Println("Waiting for response.")

	res, err := ayncResult.Get(10 * time.Second)
	if err != nil {
		panic(err)
	}

	log.Printf("result: %+v of type %+v", res, reflect.TypeOf(res))
}
