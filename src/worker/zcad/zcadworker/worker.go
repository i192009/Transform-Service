package zcadworker

import (
	"gitlab.zixel.cn/go/framework/config"
	"gitlab.zixel.cn/go/framework/logger"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var log = logger.Get()

func Start() error {

	log.Info("zcadworker start ...")
	// Create a Temporal Client to communicate with the Temporal Cluster.
	// A Temporal Client is a heavyweight object that should be created just once per process.
	c, err := client.Dial(client.Options{
		HostPort: config.GetString("temporal.address", "localhost:7233"),
	})

	if err != nil {
		log.Fatalln("Unable to create Temporal Client", err)
		return err
	}
	defer c.Close()

	w := worker.New(c, "zcad-queue", worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterWorkflow(ScheduleWorkflow)
	w.RegisterActivity(ZCAD_LoadFile)

	// Start listening to the Task Queue.
	ch := worker.InterruptCh()
	err = w.Run(ch)
	if err != nil {
		log.Fatalln("unable to start Worker", err)
		return err
	}

	return nil
}
