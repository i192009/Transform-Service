package main

// Workflow is a Temporal workflow function
// func Workflow(ctx workflow.Context, req *services.C2S_NewTaskReqT) (*services.C2S_NewTaskRpnT, error) {
// 	// Workflow logic here

// 	response := services.C2S_NewTaskRpnT{}

// 	ao := workflow.ActivityOptions{
// 		StartToCloseTimeout: time.Second,
// 		//Other Parameters for the Activity
// 	}

// 	ctx = workflow.WithActivityOptions(ctx, ao)

// 	if err := workflow.ExecuteActivity(ctx, "transformFile", req).Get(ctx, &response); err != nil {
// 		log.Println("Error in Executing the Workflow", err)
// 	}

// 	return &services.C2S_NewTaskRpnT{Code: 0, Msg: "Workflow Success", JobId: workflow.GetInfo(ctx).WorkflowExecution.ID}, nil
// }
