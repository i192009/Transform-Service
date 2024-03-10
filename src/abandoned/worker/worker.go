package main

import (
	"gitlab.zixel.cn/go/framework/logger"
)

var log = logger.Get()

// var k8s = k8smanager.NewK8sClientManager()

// func main() {
// 	fmt.Print("Started Worker")
// 	c, err := client.Dial(client.Options{
// 		HostPort: client.DefaultHostPort,
// 	})
// 	if err != nil {
// 		log.Fatalln("Unable to create client", err)
// 	}
// 	defer c.Close()

// 	w := worker.New(c, "transform-tasks", worker.Options{})

// 	w.RegisterWorkflowWithOptions(Workflow, workflow.RegisterOptions{
// 		Name: "transformFiles",
// 	})

// 	w.RegisterActivityWithOptions(NewTask, activity.RegisterOptions{
// 		Name: "transformFile",
// 	})

// 	err = w.Run(worker.InterruptCh())
// 	if err != nil {
// 		log.Fatalln("Unable to start worker", err)
// 	}
// }
