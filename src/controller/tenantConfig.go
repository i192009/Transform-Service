package controller

import (
	"context"
	"strconv"
	"transform2/models"
	"transform2/service"
	"transform2/services"
)

func SetTenantConfig(ctx context.Context, req *services.C2S_SetTanentConfigReqT) (*services.C2S_SetTanentConfigRpnT, error) {

	var rpn *services.C2S_SetTanentConfigRpnT

	//Parse the Request and Load Models that are in the Request
	pool, err := service.GetResourcePool(ctx, req.PoolId)
	if err != nil {
		rpn.Message = err.Error()
		rpn.StatusCode = "500"
		return rpn, err
	}

	jobSet, err := service.GetJobSet(ctx, req.JobSetId)
	if err != nil {
		rpn.Message = err.Error()
		rpn.StatusCode = "500"
		return rpn, err
	}

	intID := GenerateID()          // Invoke the Function to get a Unique ID
	id := strconv.Itoa(int(intID)) //Convert the integer to string to use as ID

	//Here we will let the Server know that a New Tenant Configuration is added and we'll have to Update it
	var NewTenantConfig models.TenantConfig
	NewTenantConfig.ConfigId = id
	NewTenantConfig.Pool = *pool
	NewTenantConfig.Set = *jobSet
	NewTenantConfig.AssociatedPools = req.AscPoolsId

	//Return the Response
	rpn.ConfigId = NewTenantConfig.ConfigId
	rpn.Message = "Configuration Added"
	rpn.StatusCode = "200"

	return rpn, nil
}
