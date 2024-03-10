package service

import (
	"context"
	"errors"
	"time"
	"transform2/abandoned/web/model"
	"transform2/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetEquityListByUserIDOrTenantId(userID, tenantId *string) ([]*model.Equity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	filter := bson.M{}
	if tenantId != nil && *tenantId != "" {
		filter["tenantId"] = tenantId
	} else {
		filter["userId"] = userID
	}
	var equityDB []*model.Equity
	cur, err := config.EquityCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var tmp model.Equity
		err := cur.Decode(&tmp)
		if err != nil {
			return nil, err
		}
		equityDB = append(equityDB, &tmp)
	}
	return equityDB, nil
}

func UpdateEquity(equityId string, equityDB *model.Equity) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	if equityDB == nil {
		return errors.New("equity is nil")
	}
	opts := options.Update().SetUpsert(true)
	id, err := primitive.ObjectIDFromHex(equityId)
	if err != nil {
		return err
	}
	equityDB.UpdateTime = time.Now()
	filter := bson.M{
		"_id": id,
	}
	_, err = config.EquityCollection.UpdateOne(ctx, filter, bson.M{"$set": equityDB}, opts)
	return err
}
