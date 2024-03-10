package grpcserver

import (
	"context"
	"fmt"
	"transform2/config"
	"transform2/services"

	"gitlab.zixel.cn/go/framework"
	"go.temporal.io/sdk/client"
)

// Grpc Request will execute the Workflow
func (s *TransformServer) CreateJob(ctx context.Context, req *services.C2S_CreateJobReq) (*services.S2C_CreateJobRpn, error) {
	fmt.Println("Request Came")
	var rpn services.S2C_CreateJobRpn

	var headers framework.CommonHeaders
	if err := framework.GetCommonHeaders(ctx, &headers); err != nil {
		log.Error(err.Error())
		rpn.StatusCode = 500
		rpn.Message = err.Error()
		return &rpn, err
	}

	switch req.JobType {
	case 0:
		log.Infof("ZCAD Request %v, %v", req, req.Parameters)
		workflowOptions := client.StartWorkflowOptions{
			ID:        config.RandomString(32),
			TaskQueue: "zcad-queue",
		}

		log.Debug("Starting Workflow")
		reply, err := s.WorkflowClient.ExecuteWorkflow(ctx, workflowOptions, "ScheduleWorkflow", req.StorageToken, req.Parameters)
		if err != nil {
			log.Errorf("Failed to start workflow: %v", err)
			return nil, err
		}

		var res error
		if err := reply.Get(ctx, &res); err != nil {
			log.Errorf("Failed to Get Result: %v", err)
			return nil, err
		}

		log.Infof("Result of the Request %v", rpn)
	}

	return nil, nil
}

// Get task details
func (s *TransformServer) GetJobInfo(context.Context, *services.C2S_GetJobInfoReq) (*services.S2C_GetJobInfoRpn, error) {
	return nil, nil
}

// Cancel a Task
func (s *TransformServer) CancelJob(context.Context, *services.C2S_CancelJobReq) (*services.S2C_CancelJobRpn, error) {
	return nil, nil
}

func (s *TransformServer) mustEmbedUnimplementedTransformV2Server() {}
