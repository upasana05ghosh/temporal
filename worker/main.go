package main

import (
	"log"
	"temporal/greeting"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	//client to communicate with temporal cluster
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	//greeting-tasks -> Name of the task queue, it's maintained by temporal server and polled by the worker
	// task queue and client object  is supplied to worker object
	w := worker.New(c, "greeting-tasks", worker.Options{})

	//Every workflow definition must be register with atleast one worker to proceed
	w.RegisterWorkflow(greeting.Greet)

	//Run will start the worker which will begin a long-poll on greeting-task task queue
	err = w.Run(worker.InterruptCh())

	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
