package service

import (
	"context"
	"gitlab.zixel.cn/go/framework"
	"gitlab.zixel.cn/go/framework/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"transform2/config"
	"transform2/models"
	"transform2/services"
)

var log = logger.Get()

// ValidateJobTypeById returns true if the record with the given ID does not exist.
func ValidateJobTypeById(ctx context.Context, id string) bool {
	filter := bson.M{
		"JobTypeId": id,
	}

	var result models.JobType

	if err := config.JobTypeCollection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return true // No record found, return true
		}
		log.Errorf("Error validating JobType by ID: %v", err)
		return false // Some other error occurred, return false
	}

	log.Infof("Record with ID %s already exists in the Database", id)
	return false // Record exists, return false
}

// RemoveJobTypeDuplicates function to remove duplicate job types
func RemoveJobTypeDuplicates(ctx context.Context, jobTypes []*services.JobType) []*services.JobType {
	keys := make(map[string]bool)
	uniqueJobSets := []*services.JobType{}

	for _, rp := range jobTypes {
		if _, value := keys[rp.JobTypeId]; !value {
			keys[rp.JobTypeId] = true
			uniqueJobSets = append(uniqueJobSets, rp)
		}
	}

	return uniqueJobSets
}

func GetJobSet(ctx context.Context, req string) (*models.JobSet, error) {

	var rpn models.JobSet
	filter := bson.M{"JobSetId": req}
	res := config.JobSetCollection.FindOne(ctx, filter).Decode(&rpn)
	if res == mongo.ErrNoDocuments {
		return nil, framework.NewServiceError(framework.ERR_SYS_PARSE_PARAMS, "No Such Job Set in the Database")
	}

	return &rpn, nil

}

func GetJobType(ctx context.Context, req *services.C2S_GetJobTypeReqT) (*models.JobType, error) {

	var rpn models.JobType
	filter := bson.M{"JobTypeId": req.JobTypeId}
	res := config.JobTypeCollection.FindOne(ctx, filter).Decode(&rpn)
	if res == mongo.ErrNoDocuments {
		return nil, framework.NewServiceError(framework.ERR_SYS_PARSE_PARAMS, "No Such Job Type in the Database")
	}

	return &rpn, nil

}

// AddJobType stores the JobType in the Database.
func AddJobType(ctx context.Context, jobType *models.JobType) (*string, error) {
	// Add the Record to the Database
	_, err := config.JobTypeCollection.InsertOne(ctx, &jobType)
	if err != nil {
		log.Errorf("Error adding the record to the database: %v", err)
		return nil, framework.NewServiceError(framework.ERR_SYS_DATABASE, err.Error())
	}

	// Record added to the Database
	return &jobType.JobTypeId, nil
}

// RemoveJobType removes the JobType with the given ID from the Database.
func RemoveJobType(ctx context.Context, jobTypeId string) error {
	filter := bson.M{"JobTypeId": jobTypeId}
	result, _ := config.JobTypeCollection.DeleteOne(ctx, filter)
	//log.Info("Test 2: MONGO Check %v", result)
	if result.DeletedCount == 0 {
		log.Errorf("Error removing JobType from the database")
		return framework.NewServiceError(framework.ERR_SYS_DATABASE, "No Document Found")
	} else {
		return nil
	}
}

// UpdateJobType updates the JobType in the Database.
func UpdateJobType(ctx context.Context, jobType models.JobType) error {
	// Update the Record in the Database
	filter := bson.M{"JobTypeId": jobType.JobTypeId}
	update := bson.D{
		{"$set", bson.D{
			{"ImageUrl", jobType.ImageUrl},
			{"ReScript", jobType.ReScript},
			{"ScScript", jobType.ScScript},
			{"JeScript", jobType.JeScript},
			{"FixedParameters", jobType.FixedParameters},
		}},
	}

	_, err := config.JobTypeCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("Error updating JobType in the database: %v", err)
		return framework.NewServiceError(framework.ERR_SYS_DATABASE, err.Error())
	}

	// Record Updated in the Database
	return nil
}

// AddJobSet stores the JobType in the Database.
func AddJobSet(ctx context.Context, jobSet *models.JobSet) (string, error) {
	// Add the Record to the Database
	_, err := config.JobSetCollection.InsertOne(ctx, &jobSet)
	if err != nil {
		log.Errorf("Error adding the record to the database: %v", err)
		return "", framework.NewServiceError(framework.ERR_SYS_DATABASE, err.Error())
	}

	// Record added to the Database
	return jobSet.JobSetId, nil
}

// RemoveJobSet removes the JobType with the given ID from the Database.
func RemoveJobSet(ctx context.Context, jobSetId string) error {
	filter := bson.M{"JobSetId": jobSetId}
	result, _ := config.JobSetCollection.DeleteOne(ctx, filter)
	//log.Info("Test 2: MONGO Check %v", result)
	if result.DeletedCount == 0 {
		log.Errorf("Error removing JobSet from the database")
		return framework.NewServiceError(framework.ERR_SYS_DATABASE, "No Document Found")
	} else {
		return nil
	}
}

// UpdateJobSet updates the JobSet in the Database.
func UpdateJobSet(ctx context.Context, jobSet models.JobSet) error {
	// Update the Record in the Database
	filter := bson.M{"JobSetId": jobSet.JobSetId}
	update := bson.D{
		{"$set", bson.D{
			{"JobTypeId", jobSet.JobTypeId},
			{"Name", jobSet.Name},
			{"Total", jobSet.Total},
			{"FixedParameters", jobSet.FixedParameters},
		}},
	}

	_, err := config.JobSetCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("Error updating JobSet in the database: %v", err)
		return framework.NewServiceError(framework.ERR_SYS_DATABASE, err.Error())
	}

	// Record Updated in the Database
	return nil
}

// HandleJobTypeQueries function with pagination support
func HandleJobTypeQueries(ctx context.Context, filter bson.M, skip, page, limit int32) ([]models.JobType, error) {
	var response []models.JobType

	cursor, err := config.JobTypeCollection.Find(ctx, filter)
	if err != nil {
		log.Errorf("Error querying JobTypeCollection: %v", err)
		return nil, framework.NewServiceError(framework.ERR_SYS_SERVER, "Error Querying Database")
	}
	defer cursor.Close(ctx)

	// Skip to the requested page
	for i := int32(0); i < skip; i++ {
		if cursor.Next(ctx) {
			if err := cursor.Decode(&models.JobType{}); err != nil {
				log.Errorf("Error decoding resource pool: %v", err)
				return nil, framework.NewServiceError(framework.ERR_SYS_SERVER, "Error Decoding Results")
			}
		} else {
			break
		}
	}

	// Read and append results within the pagination range
	for i := int32(0); i < limit; i++ {
		if cursor.Next(ctx) {
			var rp models.JobType
			if err := cursor.Decode(&rp); err != nil {
				log.Errorf("Error decoding resource pool: %v", err)
				return nil, framework.NewServiceError(framework.ERR_SYS_SERVER, "Error Decoding Results")
			}
			response = append(response, rp)
		} else {
			break
		}
	}

	return response, nil
}

// HandleJobSetQueries function to handle querying job sets with pagination support
func HandleJobSetQueries(ctx context.Context, filter bson.M, skip, page, limit int32) ([]models.JobSet, error) {
	var response []models.JobSet

	cursor, err := config.JobSetCollection.Find(ctx, filter)
	if err != nil {
		log.Errorf("Error querying JobSetCollection: %v", err)
		return nil, framework.NewServiceError(framework.ERR_SYS_SERVER, "Error Querying Database")
	}
	defer cursor.Close(ctx)

	// Skip to the requested page
	for i := int32(0); i < skip; i++ {
		if cursor.Next(ctx) {
			if err := cursor.Decode(&models.JobSet{}); err != nil {
				log.Errorf("Error decoding job set: %v", err)
				return nil, framework.NewServiceError(framework.ERR_SYS_SERVER, "Error Decoding Results")
			}
		} else {
			break
		}
	}

	// Read and append results within the pagination range
	for i := int32(0); i < limit; i++ {
		if cursor.Next(ctx) {
			var js models.JobSet
			if err := cursor.Decode(&js); err != nil {
				log.Errorf("Error decoding job set: %v", err)
				return nil, framework.NewServiceError(framework.ERR_SYS_SERVER, "Error Decoding Results")
			}
			response = append(response, js)
		} else {
			break
		}
	}

	return response, nil
}

// RemoveJobSetDuplicates function to remove duplicate job sets
func RemoveJobSetDuplicates(ctx context.Context, jobSets []*services.JobSet) []*services.JobSet {
	keys := make(map[string]bool)
	uniqueJobSets := []*services.JobSet{}

	for _, js := range jobSets {
		if _, value := keys[js.JobSetId]; !value {
			keys[js.JobSetId] = true
			uniqueJobSets = append(uniqueJobSets, js)
		}
	}

	return uniqueJobSets
}
