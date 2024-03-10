package controller

import (
	"context"
	"gitlab.zixel.cn/go/framework"
	"gitlab.zixel.cn/go/framework/logger"
	"strconv"
	"transform2/models"
	"transform2/service"
	"transform2/services"
)

var log = logger.Get()

func AddJobType(ctx context.Context, req *services.C2S_AddJobTypeReqT) (*services.S2C_AddJobTypeRpnT, error) {

	log.Infof("Controller Starting for Add Job Type")

	intID := GenerateID()          // Invoke the Function to get a Unique ID
	id := strconv.Itoa(int(intID)) //Convert the integer to string to use as ID

	// Parse the Request and Set the appropriate Structs to Save in the Database
	newJobType := &models.JobType{
		JobTypeId:           id,
		SystemSpecification: req.SystemSpecification,
		ImageUrl:            req.ImageUrl,
		ReScript:            req.ReScript,
		ScScript:            req.ScScript,
		JeScript:            req.JeScript,
		FixedParameters:     req.FixedParameters,
	}

	// Validate the Request before sending it to the Database
	log.Infof("Validating the Request")
	if service.ValidateJobTypeById(ctx, newJobType.JobTypeId) == false {
		log.Infof("Record Duplication, Invalid Request")
		return nil, framework.NewServiceError(framework.ERR_SYS_PARAMETER, "Record Already Exists in the Database")
	}

	// Add the JobType to the Database
	_, err := service.AddJobType(ctx, newJobType)
	if err != nil {
		return nil, err
	}

	return &services.S2C_AddJobTypeRpnT{
		JobTypeId:  id,
		Message:    "Job Type Added to the Database",
		StatusCode: "200",
	}, nil
}

func GetJobSet(ctx context.Context, req *services.C2S_GetJobSetReqT) (*services.C2S_GetJobSetRpnT, error) {

	var rpn *services.C2S_GetJobSetRpnT
	// Call the Service
	res, err := service.GetJobSet(ctx, req.JobTypeId)
	if err != nil {
		rpn.Code = 500
		return rpn, err
	}

	rpn.Jobset = &services.JobSet{
		JobSetId:        res.JobSetId,
		JobTypeId:       res.JobTypeId,
		Name:            res.Name,
		FixedParameters: res.FixedParameters,
		Total:           res.Total,
	}

	//Return the Response
	return rpn, nil

}

func GetJobType(ctx context.Context, req *services.C2S_GetJobTypeReqT) (*services.C2S_GetJobTypeRpnT, error) {

	var rpn *services.C2S_GetJobTypeRpnT
	// Call the Service
	res, err := service.GetJobType(ctx, req)
	if err != nil {
		rpn.StatusCode = "500"
		rpn.Message = err.Error()
		return rpn, err
	}

	rpn.JobType = &services.JobType{
		JobTypeId:           res.JobTypeId,
		FixedParameters:     res.FixedParameters,
		ScScript:            res.ScScript,
		SystemSpecification: res.SystemSpecification,
		ImageUrl:            res.ImageUrl,
		JeScript:            res.JeScript,
		ReScript:            res.ReScript,
	}
	rpn.StatusCode = "200"
	rpn.Message = "Job Type Found"

	//Return the Response
	return rpn, nil

}

func RemoveJobType(ctx context.Context, req *services.C2S_RemoveJobTypeReqT) (*services.S2C_RemoveJobTypeRpnT, error) {

	log.Infof("Controller Starting for Remove Job Type")

	//Remove the Requested ID from the Database
	log.Info("Calling the Service to Remove the JobType")
	err := service.RemoveJobType(ctx, req.GetJobTypeId())
	if err != nil {
		log.Info("Error Occurred :", err)
		return &services.S2C_RemoveJobTypeRpnT{StatusCode: "500", Message: err.Error()}, nil
	}
	ss := "Job Type with ID: " + req.JobTypeId + " is Deleted"
	res := &services.S2C_RemoveJobTypeRpnT{StatusCode: "200", Message: ss}

	return res, err
}

func SetJobType(ctx context.Context, req *services.C2S_SetJobTypeReqT) (*services.S2C_SetJobTypeRpnT, error) {

	log.Infof("Starting the Controller for Update Job Type")

	// Parse the Request and Set the appropriate Structs to Update in the Database
	updatedJobType := &models.JobType{
		JobTypeId:       req.JobTypeId,
		ImageUrl:        req.Image,
		ReScript:        req.ReScript,
		ScScript:        req.ScScript,
		JeScript:        req.JeScript,
		FixedParameters: req.FixedParameters,
	}

	// Validate the Request before updating it in the Database
	log.Infof("Calling the Validate Job Type By Id")
	if service.ValidateJobTypeById(ctx, updatedJobType.JobTypeId) == false {
		log.Infof("Valid Request")

		// Update the JobType in the Database
		err := service.UpdateJobType(ctx, *updatedJobType)
		if err != nil {
			log.Fatalf("Error Occured While Updating: ", err.Error())
			return &services.S2C_SetJobTypeRpnT{
				JobTypeId:  req.JobTypeId,
				Message:    "No Updates Made",
				StatusCode: "500",
			}, nil
		}

		return &services.S2C_SetJobTypeRpnT{
			JobTypeId:  req.JobTypeId,
			Message:    "Job Type Updated in the Database",
			StatusCode: "200",
		}, nil
	}

	log.Infof("Invalid Request")
	return &services.S2C_SetJobTypeRpnT{
		JobTypeId:  req.JobTypeId,
		Message:    "Record Not Found",
		StatusCode: "404",
	}, nil
}

func AddJobSet(ctx context.Context, req *services.C2S_AddJobSetReqT) (*services.C2S_AddJobSetRpnT, error) {

	log.Infof("Starting Controller for Add Job Set")
	var rpn services.C2S_AddJobSetRpnT

	intID := GenerateID()          // Invoke the Function to get a Unique ID
	id := strconv.Itoa(int(intID)) //Convert the integer to string to use as ID

	//Process the Request
	newJobSet := models.JobSet{
		JobSetId:        id,
		JobTypeId:       req.JobTypeId,
		Name:            req.Name,
		Total:           req.Total,
		FixedParameters: req.FixedParameters,
	}

	//Validate, if the Associated JobType Exists or Not
	if service.ValidateJobTypeById(ctx, newJobSet.JobTypeId) == true {
		log.Error("Associated JobType Does not Exist in the Database")
		rpn.Code = 402
		rpn.Message = framework.NewServiceError(framework.ERR_SYS_PARAMETER, "JobType Does not Exist in the Database").Error()
		return &rpn, nil
	}

	//Associated Job Type Exists, Proceed with the Request
	if _, err := service.AddJobSet(ctx, &newJobSet); err != nil {
		log.Error("Error Adding the JobSet to the Database")
		rpn.Code = 402
		rpn.Message = err.Error()
		return &rpn, nil
	}

	rpn.JobSetId = newJobSet.JobSetId
	rpn.Code = 200
	rpn.Message = "Job Set Added Successfully"

	log.Infof("Service Finished")

	return &rpn, nil
}

func RemoveJobSet(ctx context.Context, req *services.C2S_RemoveJobSetReqT) (*services.C2S_RemoveJobSetRpnT, error) {
	log.Infof("Controller for Remove Job Set")
	var rpn services.C2S_RemoveJobSetRpnT

	//Process the Request
	if err := service.RemoveJobSet(ctx, req.JobSetId); err != nil {
		log.Error(err.Error())
		rpn.Code = 500
		rpn.Message = err.Error()
		return &rpn, nil
	}

	//Return the Response
	rpn.Code = 200
	rpn.Message = "Job Set Deleted Successfully"
	return &rpn, nil
}

func SetJobSet(ctx context.Context, req *services.C2S_SetJobSetReqT) (*services.C2S_SetJobSetRpnT, error) {
	log.Infof("Controller for Update Job Set")
	var rpn services.C2S_SetJobSetRpnT

	//Parse the Request
	jobSet := models.JobSet{
		JobSetId:        req.JobSetId,
		JobTypeId:       req.JobTypeId,
		FixedParameters: req.FixedParameters,
		Total:           req.Total,
		Name:            req.Name,
	}

	//Validate the Updated Job Set before passing the request

	if service.ValidateJobTypeById(ctx, jobSet.JobTypeId) == true {
		log.Error("No Associated Job Type Present in the Database")
		rpn.Code = 402
		rpn.Message = framework.NewServiceError(framework.ERR_SYS_PARAMETER, "No Associated Job Type Present in the Database").Error()
		return &rpn, nil
	}

	if err := service.UpdateJobSet(ctx, jobSet); err != nil {
		log.Error(err.Error())
		rpn.Code = 500
		rpn.Message = err.Error()
		return &rpn, nil
	}

	rpn.Code = 200
	rpn.Message = "Job Set Updated Successfully"
	rpn.JobSetId = jobSet.JobSetId

	// Return the Response
	return &rpn, nil
}

// QueryJobSet function to handle querying job sets with filters and pagination
func QueryJobSet(ctx context.Context, req *services.C2S_QueryJobSetReqT) (*services.C2S_QueryJobSetRpnT, error) {
	rpn := &services.C2S_QueryJobSetRpnT{}

	// Build filters based on the request
	filters := models.JobSetFilter{
		JobSetIdFilter:        req.JobSetIdFilter,
		JobTypeIdFilter:       req.JobTypeIdFilter,
		NameFilter:            req.NameFilter,
		TotalFilter:           req.TotalFilter,
		FixedParametersFilter: req.FixedParameterFilter,
	}

	// Set up pagination parameters
	skip, page, limit := req.Skip, req.Page, req.Limit

	// Execute queries for each filter
	for i := 0; i < 5; i++ {
		// Create a map for the current filter
		dbFilter := make(map[string]interface{})

		// Set the current filter based on the iteration
		switch i {
		case 0:
			if filters.JobSetIdFilter != "" {
				dbFilter["JobSetId"] = filters.JobSetIdFilter
			} else {
				continue
			}
		case 1:
			if filters.JobTypeIdFilter != "" {
				dbFilter["JobTypeId"] = filters.JobTypeIdFilter
			} else {
				continue
			}
		case 2:
			if filters.NameFilter != "" {
				dbFilter["Name"] = filters.NameFilter
			} else {
				continue
			}
		case 3:
			dbFilter["Total"] = filters.TotalFilter
		case 4:
			if filters.FixedParametersFilter != "" {
				dbFilter["FixedParameters"] = filters.FixedParametersFilter
			} else {
				continue
			}
		}

		// Execute the query using the controller function
		result, err := service.HandleJobSetQueries(ctx, dbFilter, skip, page, limit)
		if err != nil {
			return nil, err
		}

		// Append the results to the response
		for _, sets := range result {
			newJobSet := &services.JobSet{
				JobSetId:        sets.JobSetId,
				JobTypeId:       sets.JobTypeId,
				Name:            sets.Name,
				Total:           sets.Total,
				FixedParameters: sets.FixedParameters,
			}

			rpn.Jobsets = append(rpn.Jobsets, newJobSet)
		}
	}

	// Set success code and message
	rpn.Code = 200

	// Call the Service to remove the Duplicates
	service.RemoveJobSetDuplicates(ctx, rpn.Jobsets)
	return rpn, nil
}

func QueryJobType(ctx context.Context, req *services.C2S_QueryJobTypeReqT) (*services.C2S_QueryJobTypeRpnT, error) {

	// Initialize rpn
	rpn := &services.C2S_QueryJobTypeRpnT{}

	// Build filters based on the request
	filters := models.JobTypeFilter{
		JobTypeIdFilter:           req.JobTypeIdFilter,
		SystemSpecificationFilter: req.SystemSpecificationFilter,
		ImageUrlFilter:            req.ImageUrlFilter,
		ReScriptFilter:            req.ReScriptFilter,
		ScScriptFilter:            req.ScScriptFilter,
		JeScriptFilter:            req.JeScriptFilter,
		FixedParametersFilter:     req.FixedParametersFilter,
	}

	// Set up pagination parameters
	var skip, page, limit int32
	skip = req.Skip
	page = req.Page
	limit = req.Limit

	// Execute queries for each filter
	for i := 0; i < 7; i++ {
		// Create a map for the current filter
		dbFilter := make(map[string]interface{})

		// Set the current filter based on the iteration
		switch i {
		case 0:
			if filters.JobTypeIdFilter != "" {
				dbFilter["JobTypeId"] = filters.JobTypeIdFilter
			} else {
				continue
			}
		case 1:
			if filters.SystemSpecificationFilter != 0 {
				dbFilter["SystemSpecification"] = filters.SystemSpecificationFilter
			} else {
				continue
			}
		case 2:
			if filters.ImageUrlFilter != "" {
				dbFilter["ImageUrl"] = filters.ImageUrlFilter
			} else {
				continue
			}
		case 3:
			if filters.FixedParametersFilter != nil {
				dbFilter["FixedParameters"] = filters.FixedParametersFilter
			} else {
				continue
			}
		case 4:
			if filters.ReScriptFilter != "" {
				dbFilter["ReScript"] = filters.ReScriptFilter
			} else {
				continue
			}
		case 5:
			if filters.ScScriptFilter != "" {
				dbFilter["ScScript"] = filters.ScScriptFilter
			} else {
				continue
			}
		case 6:
			if filters.JeScriptFilter != "" {
				dbFilter["JeScript"] = filters.JeScriptFilter
			} else {
				continue
			}
		}

		// Execute the query using the controller function
		if result, err := service.HandleJobTypeQueries(ctx, dbFilter, skip, page, limit); err != nil {
			return nil, err
		} else {
			// Apply pagination logic
			startIndex := int(page-1)*int(limit) + int(skip)
			endIndex := startIndex + int(limit)
			if startIndex < 0 {
				startIndex = 0
			}
			if endIndex > len(result) {
				endIndex = len(result)
			}

			// Append the results to the response within the pagination range
			for _, sets := range result[startIndex:endIndex] {
				newPool := &services.JobType{
					JobTypeId:           sets.JobTypeId,
					SystemSpecification: sets.SystemSpecification,
					ImageUrl:            sets.ImageUrl,
					FixedParameters:     sets.FixedParameters,
					ReScript:            sets.ReScript,
					ScScript:            sets.ScScript,
					JeScript:            sets.JeScript,
				}

				rpn.JobTypes = append(rpn.JobTypes, newPool)
			}
		}
	}

	// Set success code and message
	rpn.StatusCode = "200"

	// Call the Service to remove the Duplicates
	service.RemoveJobTypeDuplicates(ctx, rpn.JobTypes)
	return rpn, nil
}
