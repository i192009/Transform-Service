package controller

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gitlab.zixel.cn/go/framework"
	"strconv"
	"transform2/models"
	"transform2/service"
	"transform2/services"
)

func AddResourcePool(ctx context.Context, req *services.C2S_AddResourcePoolReqT) (*services.C2S_AddResourcePoolRpnT, error) {

	fmt.Println("Controller for Add Resource Pool")

	intID := GenerateID()          // Invoke the Function to get a Unique ID
	id := strconv.Itoa(int(intID)) //Convert the integer to string to use as ID

	// Parse the Request and Set the appropriate Structs to Save in the Database
	newResourcePool := &models.ResourcePool{
		ResourcePoolID: id,
		Name:           req.Name,
		NameSpace:      req.NameSpace,
		IsShared:       req.IsShared,
		ScalingLimit:   req.ScalingLimit,
		QueueLimit:     req.QueueLimit,
		Fixed:          req.Fixed,
		DefaultTaskSet: req.DefaultTaskSet,
	}

	// Set ResourceLimits for each job type
	for _, rl := range req.ResourceLimit {
		jobTypeId := rl.JobTypeId
		scalingLimit := rl.ScalingLimit

		newResourcePool.ResourceLimits[jobTypeId] = &models.ResourceLimitOfTask{
			JobTypeId:    jobTypeId,
			ScalingLimit: scalingLimit,
		}
	}

	// Validate the Request before sending it to the Database
	fmt.Println("Calling the Validate Resource Pool By Id")
	if service.ValidateRpById(ctx, newResourcePool.ResourcePoolID) == false {
		log.Infof("Record Duplication, Invalid Request")
		return nil, framework.NewServiceError(framework.ERR_SYS_PARAMETER, "Record Already Exists in the Database")
	}

	// Add the Resource Pool to the Database
	_, err := service.AddResourcePool(ctx, *newResourcePool)
	if err != nil {
		return nil, err
	}

	return &services.C2S_AddResourcePoolRpnT{
		ResourcePoolId: id,
		NameSpace:      newResourcePool.NameSpace,
		Message:        "Resource Pool Added to the Database",
		StatusCode:     "200",
	}, nil
}

func RemoveResourcePool(ctx context.Context, req *services.C2S_RemoveResourcePoolReqT) (*services.C2S_RemoveResourcePoolRpnT, error) {

	log.Infof("Controller for Remove Resource Pool")

	if err := service.RemoveResourcePool(ctx, req.PoolId); err != nil {
		return &services.C2S_RemoveResourcePoolRpnT{StatusCode: "404", Message: err.Error()}, nil
	}
	ss := "Resource Pool with ID: " + req.PoolId + " is Deleted"
	res := &services.C2S_RemoveResourcePoolRpnT{StatusCode: "200", Message: ss}

	return res, nil
}

func SetResourcePool(ctx context.Context, req *services.C2S_SetResourcePoolReqT) (*services.C2S_SetResourcePoolRpnT, error) {

	log.Infof("Controller for Update Resource Pool")

	// Parse the Request and Set the appropriate Structs to Save in the Database
	updateRP := &models.ResourcePool{
		ResourcePoolID:  req.ResourcePoolId, // Use the existing ID for updating
		Name:            req.Name,
		NameSpace:       req.NameSpace,
		IsShared:        req.IsShared,
		ScalingStrategy: req.ScalingStrategy,
		ScalingLimit:    req.ScalingLimit,
		QueueLimit:      req.QueueLimit,
		Fixed:           req.Fixed,
		DefaultTaskSet:  req.DefaultTaskSet,
	}

	// Set ResourceLimit
	for _, rl := range req.ResourceLimit {

		jobTypeId := rl.JobTypeId
		scalingLimit := rl.ScalingLimit

		updateRP.ResourceLimits[jobTypeId] = &models.ResourceLimitOfTask{
			//JobTypeId:    rl.JobTypeId,
			ScalingLimit: scalingLimit,
		}
	}

	// Validate the Request before updating it in the Database
	fmt.Println("Calling the Validate Resource Pool By Id")
	if service.ValidateRpById(ctx, updateRP.ResourcePoolID) == false {
		log.Infof("Valid Request")

		// Update the ResourcePool in the Database
		err := service.UpdateResourcePool(ctx, *updateRP)
		if err != nil {
			return nil, err
		}

		return &services.C2S_SetResourcePoolRpnT{
			ResourcePoolId: req.ResourcePoolId,
			Message:        "Resource Pool Updated in the Database",
			StatusCode:     "200",
		}, nil
	}

	log.Infof("Invalid Request")
	return &services.C2S_SetResourcePoolRpnT{
		Message:    "Record Does Not Exist",
		StatusCode: "404",
	}, nil
}

func GetResourcePool(ctx context.Context, req *services.C2S_GetResourcePoolReqT) (*services.C2S_GetResourcePoolRpnT, error) {
	log.Infof("Controller for Get Resource Pool")

	//Initialize the Response parameter
	rpn := &services.C2S_GetResourcePoolRpnT{}
	rpn = new(services.C2S_GetResourcePoolRpnT)

	//Parse the request
	var resPool models.ResourcePool
	resPool.ResourcePoolID = req.PoolId
	//resPool.Name = *req.PoolName

	//Pass the request to the Controller to see if the pool exists in the database
	if err := service.ValidateRpById(ctx, resPool.ResourcePoolID); err == true {
		log.Infof("No Record Present in the Database")
		rpn.Pool = nil
		rpn.Code = proto.String("200")
		rpn.Message = proto.String("No Record Present in the Database")
		return rpn, nil
	}

	//Return if error occurs
	res, err := service.GetResourcePool(ctx, resPool.ResourcePoolID)
	if err != nil {
		rpn.Code = proto.String("200")
		rpn.Message = proto.String(err.Error())
		return rpn, nil
	}

	log.Infof("Document Found")

	//Set the Response Pool
	pool := &services.ResourcePool{
		ResourcePoolId:  res.ResourcePoolID,
		Name:            res.Name,
		NameSpace:       res.NameSpace,
		IsShared:        res.IsShared,
		ScalingLimit:    res.ScalingLimit,
		ScalingStrategy: res.ScalingStrategy,
		QueueLimit:      res.QueueLimit,
		Fixed:           res.Fixed,
		DefaultTaskSet:  res.DefaultTaskSet,
	}

	// Set ResourceLimit
	for _, rl := range res.ResourceLimits {

		jobTypeId := rl.JobTypeId
		scalingLimit := rl.ScalingLimit

		pool.ResourceLimits[jobTypeId] = &services.ResourceLimitOfTask{
			ScalingLimit: scalingLimit,
		}
	}

	log.Infof("Return the Response")

	response := services.C2S_GetResourcePoolRpnT{
		Pool: pool,
	}

	return &response, nil
}

// QueryResourcePool function to handle querying resource pools with filters and pagination
func QueryResourcePool(ctx context.Context, req *services.C2S_QueryResourcePoolReqT) (*services.C2S_QueryResourcePoolRpnT, error) {
	rpn := &services.C2S_QueryResourcePoolRpnT{}

	// Build filters based on the request
	filters := models.ResourcePoolFilter{
		ResourcePoolIDFilter:  req.ResourcePoolIdFilter,
		NameSpaceFilter:       req.NameSpaceFilter,
		NameFilter:            req.NameFilter,
		IsSharedFilter:        req.IsSharedFilter,
		ScalingStrategyFilter: req.ScalingStrategyFilter,
		ScalingLimitFilter:    req.ScalingLimitFilter,
		CustomFilter:          req.CustomFilter,
	}

	// Set up pagination parameters
	skip, page, limit := req.Skip, req.Page, req.Limit

	// Execute queries for each filter
	for i := 0; i < 6; i++ {
		// Create a map for the current filter
		dbFilter := make(map[string]interface{})

		// Set the current filter based on the iteration
		switch i {
		case 0:
			if filters.ResourcePoolIDFilter != "" {
				dbFilter["ResourcePoolId"] = filters.ResourcePoolIDFilter
			} else {
				continue
			}
		case 1:
			if filters.NameFilter != "" {
				dbFilter["Name"] = filters.NameFilter
			} else {
				continue
			}
		case 2:
			dbFilter["IsShared"] = filters.IsSharedFilter
		case 3:
			if filters.ScalingStrategyFilter != 0 {
				dbFilter["ScalingStrategy"] = filters.ScalingStrategyFilter
			} else {
				continue
			}
		case 4:
			if filters.ScalingLimitFilter != 0 {
				dbFilter["ScalingLimit"] = filters.ScalingLimitFilter
			} else {
				continue
			}
		case 5:
			if filters.NameSpaceFilter != "" {
				dbFilter["NameSpace"] = filters.NameSpaceFilter
			} else {
				continue
			}
		}

		// Execute the query using the controller function
		result, err := service.HandleQueries(ctx, dbFilter, skip, page, limit)
		if err != nil {
			return rpn, nil
		}

		// Append the results to the response
		for _, pool := range result {
			newPool := &services.ResourcePool{
				ResourcePoolId:  pool.ResourcePoolID,
				Name:            pool.Name,
				NameSpace:       pool.NameSpace,
				IsShared:        pool.IsShared,
				ScalingLimit:    pool.ScalingLimit,
				ScalingStrategy: pool.ScalingStrategy,
				QueueLimit:      pool.QueueLimit,
				Fixed:           pool.Fixed,
				DefaultTaskSet:  pool.DefaultTaskSet,
			}

			for _, rl := range pool.ResourceLimits {

				typeId := rl.JobTypeId
				scaling := rl.ScalingLimit

				newPool.ResourceLimits[typeId] = &services.ResourceLimitOfTask{
					ScalingLimit: scaling,
				}
			}
			rpn.ResourcePools = append(rpn.ResourcePools, newPool)
		}
	}

	// Set success code and message
	rpn.Code = proto.String("200")
	rpn.Message = proto.String("Resource Pools retrieved successfully")

	// Call the Service to remove the Duplicates
	rpn.ResourcePools = service.RemoveDuplicates(ctx, rpn.ResourcePools)

	return rpn, nil
}
