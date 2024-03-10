package service

import (
	"context"
	"gitlab.zixel.cn/go/framework"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"transform2/config"
	"transform2/models"
	"transform2/services"
)

// returns Resource Pool if it exists in the Database.
func GetResourcePool(ctx context.Context, id string) (*models.ResourcePool, error) {
	filter := bson.M{
		"ResourcePoolId": id,
	}

	//Create a Response Parameter
	var result models.ResourcePool

	if err := config.RpTypeCollection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, framework.NewServiceError(framework.ERR_SYS_PARAMETER, "No Resource Pool Exists in the Database") // No record found, return true
		}
		log.Errorf("Error validating ResourcePool by ID")
		return nil, framework.NewServiceError(framework.ERR_SYS_SERVER, err.Error()) // Some other error occurred, return false
	}

	log.Infof("Record with ID %s Found", id)
	return &result, nil // Record exists, return false
}

// ValidateRpById returns true if the record with the given ID does not exist.
func ValidateRpById(ctx context.Context, id string) bool {
	filter := bson.M{
		"ResourcePoolId": id,
	}

	var result models.ResourcePool

	if err := config.RpTypeCollection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Errorf("No record in the Database for Request")
			return true // No record found, return true
		}
		log.Errorf("Error validating ResourcePool by ID")
		return false // Some other error occurred, return false
	}

	log.Infof("Record with ID %s already exists in the Database", id)
	return false // Record exists, return false
}

// AddResourcePool stores the ResourcePool in the Database.
func AddResourcePool(ctx context.Context, rp models.ResourcePool) (*string, error) {
	// Add the Record to the Database
	_, err := config.RpTypeCollection.InsertOne(ctx, &rp)
	if err != nil {
		log.Errorf("Error adding the record to the database: %v", err)
		return nil, framework.NewServiceError(framework.ERR_SYS_DATABASE, err.Error())
	}

	// Record added to the Database
	return &rp.ResourcePoolID, nil
}

// RemoveResourcePool removes the ResourcePool with the given ID from the Database.
func RemoveResourcePool(ctx context.Context, rpId string) error {
	filter := bson.M{"ResourcePoolId": rpId}
	res, _ := config.RpTypeCollection.DeleteOne(ctx, filter)
	if res.DeletedCount == 0 {
		return framework.NewServiceError(framework.ERR_SYS_DATABASE, "No Record Found")
	}
	return nil
}

// UpdateResourcePool updates the ResourcePool in the Database.
func UpdateResourcePool(ctx context.Context, rp models.ResourcePool) error {
	// Update the Record in the Database
	filter := bson.M{"ResourcePoolId": rp.ResourcePoolID}
	update := bson.D{
		{"$set", bson.D{
			{"Name", rp.Name},
			{"IsShared", rp.IsShared},
			{"ScalingStrategy", rp.ScalingStrategy},
			{"ScalingLimit", rp.ScalingLimit},
			{"QueueLimit", rp.QueueLimit},
			{"Fixed", rp.Fixed},
			{"DefaultTaskSet", rp.DefaultTaskSet},
			{"ResourceLimit", rp.ResourceLimits},
		}},
	}

	_, err := config.RpTypeCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("Error updating ResourcePool in the database: %v", err)
		return framework.NewServiceError(framework.ERR_SYS_DATABASE, err.Error())
	}

	// Record Updated in the Database
	return nil
}

// HandleQueries function to handle querying resource pools with pagination support
func HandleQueries(ctx context.Context, filter bson.M, skip, page, limit int32) ([]models.ResourcePool, error) {
	var response []models.ResourcePool

	cursor, err := config.RpTypeCollection.Find(ctx, filter)
	if err != nil {
		log.Errorf("Error querying RpTypeCollection: %v", err)
		return nil, framework.NewServiceError(framework.ERR_SYS_SERVER, "Error Querying Database")
	}
	defer cursor.Close(ctx)

	// Skip to the requested page
	for i := int32(0); i < skip; i++ {
		if cursor.Next(ctx) {
			if err := cursor.Decode(&models.ResourcePool{}); err != nil {
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
			var rp models.ResourcePool
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

// RemoveDuplicates function to remove duplicate resource pools
func RemoveDuplicates(ctx context.Context, resourcePools []*services.ResourcePool) []*services.ResourcePool {
	keys := make(map[string]bool)
	uniqueResourcePools := []*services.ResourcePool{}

	for _, rp := range resourcePools {
		if _, value := keys[rp.ResourcePoolId]; !value {
			keys[rp.ResourcePoolId] = true
			uniqueResourcePools = append(uniqueResourcePools, rp)
		}
	}

	return uniqueResourcePools
}
