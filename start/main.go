package main

import (
	"context"
	"log"
	"os"
	"temporal/greeting"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "my-first-workflow",
		TaskQueue: "greeting-tasks",
	}

	//ExecuteWorkflow - returns future object since it will get response only when it's picked by worker
	we, err := c.ExecuteWorkflow(context.Background(), options, greeting.Greet, os.Args[1])
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	var result string
	//Get - block until workflow execution is finished. It will either get error or response
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}
