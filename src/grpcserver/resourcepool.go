package grpcserver

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"transform2/controller"
	"transform2/services"
)

type ResourcePoolServer struct {
	services.UnimplementedResourcePoolManagementServer
}

func (s *ResourcePoolServer) AddResourcePool(ctx context.Context, req *services.C2S_AddResourcePoolReqT) (*services.C2S_AddResourcePoolRpnT, error) {

	fmt.Println("Request Came for Add Resource Pool")

	//Pass the Request to the Controller

	response, err := controller.AddResourcePool(ctx, req)
	if err != nil {
		response.StatusCode = "500"
		response.Message = err.Error()
		return response, nil
	}

	//Return the Response
	return response, nil

}

func (s *ResourcePoolServer) RemoveResourcePool(ctx context.Context, req *services.C2S_RemoveResourcePoolReqT) (*services.C2S_RemoveResourcePoolRpnT, error) {

	log.Infof("Request Came for Remove Resource Pool")

	//Pass the Request to the Controller

	response, err := controller.RemoveResourcePool(ctx, req)
	if err != nil {
		response.StatusCode = "500"
		response.Message = err.Error()
		return response, nil
	}

	//Return the response
	return response, nil
}

func (s *ResourcePoolServer) SetResourcePool(ctx context.Context, req *services.C2S_SetResourcePoolReqT) (*services.C2S_SetResourcePoolRpnT, error) {

	log.Infof("Request Came for Update Resource Pool")

	//Pass the Request to the Controller
	response, err := controller.SetResourcePool(ctx, req)
	if err != nil {
		response.StatusCode = "500"
		response.Message = err.Error()
		return response, nil
	}

	//Return the Response

	return response, nil
}

func (s *ResourcePoolServer) GetResourcePool(ctx context.Context, req *services.C2S_GetResourcePoolReqT) (*services.C2S_GetResourcePoolRpnT, error) {
	log.Infof("Request Came for Get Resource Pool")

	var response *services.C2S_GetResourcePoolRpnT

	//Pass the Request to the Controller
	result, err := controller.GetResourcePool(ctx, req)
	if err != nil {
		response.Code = proto.String("500")
		response.Message = proto.String(err.Error())
		return response, nil
	}

	//Return the Response
	return result, nil
}

func (s *ResourcePoolServer) QueryResourcePool(ctx context.Context, req *services.C2S_QueryResourcePoolReqT) (*services.C2S_QueryResourcePoolRpnT, error) {
	log.Infof("Request Came for QueryResourcePool")

	//Pass the Request to the Controller
	rpn, err := controller.QueryResourcePool(ctx, req)
	if err != nil {
		rpn.Code = proto.String("500")
		rpn.Message = proto.String(err.Error())
		return rpn, nil
	}

	//Return The Response

	return rpn, nil
}
