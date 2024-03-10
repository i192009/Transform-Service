package grpcserver

import (
	"context"
	"fmt"
	"gitlab.zixel.cn/go/framework"
	"transform2/controller"
	"transform2/services"
)

type JobManageServer struct {
	services.UnimplementedJobManagementServer
}

func (s *JobManageServer) AddJobType(ctx context.Context, req *services.C2S_AddJobTypeReqT) (*services.S2C_AddJobTypeRpnT, error) {

	fmt.Println("Request Came for Add Job Type")
	var rpn services.S2C_AddJobTypeRpnT

	var headers framework.CommonHeaders
	if err := framework.GetCommonHeaders(ctx, &headers); err != nil {
		log.Error(err.Error())
		rpn.JobTypeId = ""
		rpn.StatusCode = "500"
		rpn.Message = err.Error()
		return &rpn, err
	}

	//Passing the Request to the Controller Layer with the Context
	response, err := controller.AddJobType(ctx, req)
	if err != nil {
		response.JobTypeId = ""
		response.StatusCode = "500"
		response.Message = err.Error()
		return response, nil
	}
	//Return the received response from the controller
	return response, nil
}

func (s *JobManageServer) RemoveJobType(ctx context.Context, req *services.C2S_RemoveJobTypeReqT) (*services.S2C_RemoveJobTypeRpnT, error) {

	log.Infof("Request Came for Remove Job Type")
	var rpn services.S2C_RemoveJobTypeRpnT

	var headers framework.CommonHeaders
	if err := framework.GetCommonHeaders(ctx, &headers); err != nil {
		log.Error(err.Error())
		rpn.StatusCode = "500"
		rpn.Message = err.Error()
		return &rpn, err
	}

	log.Info("Passing the Request to Controller")

	response, err := controller.RemoveJobType(ctx, req)
	if err != nil {
		log.Info("Error Occurred :", err)
		response.StatusCode = "500"
		response.Message = err.Error()
		return response, nil
	}

	//Return the Response
	return response, nil
}

func (s *JobManageServer) SetJobType(ctx context.Context, req *services.C2S_SetJobTypeReqT) (*services.S2C_SetJobTypeRpnT, error) {

	log.Infof("Request Came for Update Job Type")
	var rpn services.S2C_SetJobTypeRpnT

	var headers framework.CommonHeaders
	if err := framework.GetCommonHeaders(ctx, &headers); err != nil {
		log.Error(err.Error())
		rpn.StatusCode = "500"
		rpn.Message = err.Error()
		return &rpn, err
	}

	//Pass the Request to Controller
	response, err := controller.SetJobType(ctx, req)
	if err != nil {
		response.StatusCode = "500"
		response.Message = err.Error()
		return response, nil
	}

	//Return the Response
	return response, nil
}

func (s *JobManageServer) AddJobSet(ctx context.Context, req *services.C2S_AddJobSetReqT) (*services.C2S_AddJobSetRpnT, error) {
	log.Infof("Request Came for Add Job Set")
	var rpn services.C2S_AddJobSetRpnT

	var headers framework.CommonHeaders
	if err := framework.GetCommonHeaders(ctx, &headers); err != nil {
		log.Error(err.Error())
		rpn.Code = 500
		rpn.Message = err.Error()
		return &rpn, nil
	}

	//Call the Controller by passing the Request
	response, err := controller.AddJobSet(ctx, req)
	if err != nil {
		response.Code = 500
		response.Message = err.Error()
		return response, nil
	}

	//Return the Response
	return response, nil
}

func (s *JobManageServer) QueryJobSet(ctx context.Context, req *services.C2S_QueryJobSetReqT) (*services.C2S_QueryJobSetRpnT, error) {
	log.Infof("Request Came for QueryJobSet")

	//Call the Controller
	response, err := controller.QueryJobSet(ctx, req)
	if err != nil {
		response.Code = 500
		return response, nil
	}
	//Return the Response
	return response, nil
}

func (s *JobManageServer) GetJobSet(ctx context.Context, req *services.C2S_GetJobSetReqT) (*services.C2S_GetJobSetRpnT, error) {
	log.Infof("Request Came for QueryJobSet")

	//Call the Controller
	response, err := controller.GetJobSet(ctx, req)
	if err != nil {
		response.Code = 500
		return response, nil
	}
	//Return the Response
	return response, nil
}

func (s *JobManageServer) QueryJobType(ctx context.Context, req *services.C2S_QueryJobTypeReqT) (*services.C2S_QueryJobTypeRpnT, error) {
	log.Infof("Request Came for QueryJobType")

	//Call the Controller
	response, err := controller.QueryJobType(ctx, req)
	if err != nil {
		response.StatusCode = "500"
		return response, nil
	}
	//Return the Response
	return response, nil
}

func (s *JobManageServer) GetJobType(ctx context.Context, req *services.C2S_GetJobTypeReqT) (*services.C2S_GetJobTypeRpnT, error) {
	log.Infof("Request Came for GetJobType")

	//Call the Controller
	response, err := controller.GetJobType(ctx, req)
	if err != nil {
		response.StatusCode = "500"
		response.Message = err.Error()
		return response, nil
	}
	//Return the Response
	return response, nil
}

func (s *JobManageServer) RemoveJobSet(ctx context.Context, req *services.C2S_RemoveJobSetReqT) (*services.C2S_RemoveJobSetRpnT, error) {
	log.Infof("Request Came for Remove Job Set")

	//Call the Controller
	response, err := controller.RemoveJobSet(ctx, req)
	if err != nil {
		response.Code = 500
		response.Message = err.Error()
		return response, nil
	}

	//Return the Response
	return response, nil
}

func (s *JobManageServer) SetJobSet(ctx context.Context, req *services.C2S_SetJobSetReqT) (*services.C2S_SetJobSetRpnT, error) {
	log.Infof("Request Came for Update Job Set")
	var rpn services.C2S_SetJobSetRpnT
	var headers framework.CommonHeaders
	if err := framework.GetCommonHeaders(ctx, &headers); err != nil {
		log.Error(err.Error())
		rpn.Code = 500
		rpn.Message = err.Error()
		return &rpn, nil
	}

	//Call the Controller
	response, err := controller.SetJobSet(ctx, req)
	if err != nil {
		response.Code = 500
		response.Message = err.Error()
		return response, nil
	}

	//Return the Response
	return response, nil
}
