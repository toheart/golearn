package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/toheart/golearn/tTemporal/helloworld"
	"go.temporal.io/sdk/client"
	"time"
)

/**
@file:
@author: levi.Tang
@time: 2024/5/21 15:00
@description:
**/

func main() {

	// Create the client object just once per process
	c, err := client.Dial(client.Options{
		HostPort: "175.178.49.104:7233",
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        uuid.NewString(),
		TaskQueue: helloworld.GreetingTaskQueue,
	}
	// Start the Workflow
	name := "World_tly"
	log.Infof("before call workerflow %s", name)
	we, err := c.ExecuteWorkflow(context.Background(), options, helloworld.GreetingWorkflow, name)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}
	log.Infof("mid call workflow %s", name)
	// Get the results
	var greeting string
	err = we.Get(context.Background(), &greeting)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}
	log.Infof("after call workflow %s", name)
	printResults(greeting, we.GetID(), we.GetRunID())
	time.Sleep(1 * time.Second)
}

func printResults(greeting string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s\n\n", greeting)
}
