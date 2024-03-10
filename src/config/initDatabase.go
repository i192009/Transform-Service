package config

import (
	"errors"
	"gitlab.zixel.cn/go/framework/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var JobsCollection *mongo.Collection = nil
var ScriptCollection *mongo.Collection = nil
var ClientRegisterCollection *mongo.Collection = nil
var EquityCollection *mongo.Collection = nil
var JobTypeCollection *mongo.Collection = nil
var JobSetCollection *mongo.Collection = nil
var RpTypeCollection *mongo.Collection = nil

func InitMongoDB() (err error) {
	if JobsCollection = database.GetCollection("jobs"); JobsCollection == nil {
		err = errors.New("jobs collection not found")
		return
	}

	if JobSetCollection = database.GetCollection("jobSets"); JobsCollection == nil {
		err = errors.New("jobs collection not found")
		return
	}
	if ClientRegisterCollection = database.GetCollection("clientRegister"); ClientRegisterCollection == nil {
		err = errors.New("clientRegister collection not found")
		return
	}
	if EquityCollection = database.GetCollection("equity"); ClientRegisterCollection == nil {
		err = errors.New("clientRegister collection not found")
		return
	}
	if ScriptCollection = database.GetCollection("script"); JobsCollection == nil {
		err = errors.New("script collection not found")
		return
	}

	if JobTypeCollection = database.GetCollection("jobsType"); JobsCollection == nil {
		err = errors.New("jobsType collection not found")
		return
	}

	if RpTypeCollection = database.GetCollection("resourcePool"); JobsCollection == nil {
		err = errors.New("jobsType collection not found")
		return
	}
	return
}
