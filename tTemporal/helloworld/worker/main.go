package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/toheart/golearn/tTemporal/helloworld"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

/**
@file:
@author: levi.Tang
@time: 2024/5/21 14:53
@description:
**/

func main() {
	c, err := client.Dial(client.Options{
		HostPort: "175.178.49.104:7233",
	})
	if err != nil {
		log.Fatalln("unable to create temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, helloworld.GreetingTaskQueue, worker.Options{})
	w.RegisterWorkflow(helloworld.GreetingWorkflow)
	w.RegisterActivity(helloworld.GetNameAsync)
	w.RegisterActivity(helloworld.SayHello)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start worker", err)
	}
}
