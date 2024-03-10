package grpcserver

import (
	"context"
	"transform2/controller"
	"transform2/services"
)

type TenantConfigServer struct {
	services.UnimplementedTenantManagementServer
}

func (s *TenantConfigServer) SetTenantConfig(ctx context.Context, req *services.C2S_SetTanentConfigReqT) (*services.C2S_SetTanentConfigRpnT, error) {

	log.Infof("Request for Set Tenant Config Came")
	var rpn *services.C2S_SetTanentConfigRpnT
	//Call the Controller and Return the response

	res, err := controller.SetTenantConfig(ctx, req)
	if err != nil {
		rpn.Message = err.Error()
		rpn.StatusCode = "500"
		return rpn, nil
	}

	return res, nil
}

func (s *TenantConfigServer) SetDefaultTenantConfig(ctx context.Context, req *services.C2S_SetTanentConfigReqT) (*services.C2S_SetTanentConfigRpnT, error) {

	log.Infof("Request for Set Tenant Config Came")
	var rpn *services.C2S_SetTanentConfigRpnT
	//Call the Controller and Return the response

	res, err := controller.SetTenantConfig(ctx, req)
	if err != nil {
		rpn.Message = err.Error()
		rpn.StatusCode = "500"
		return rpn, nil
	}

	return res, nil
}
