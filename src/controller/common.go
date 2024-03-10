package controller

import "gitlab.zixel.cn/go/framework/database"

var (
	GenerateID, _ = database.Mongo_NewIdAllocator("generateID") //Used to generate the ID's of the new JobTypes, JobSets etc.
)
