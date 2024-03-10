package models

// JobType struct for different JobTypes
type JobType struct {
	JobTypeId           string   `json:"JobTypeId" bson:"JobTypeId"`                               //Unique Identifier for the JobType
	SystemSpecification int32    `json:"SystemSpecification,omitempty" bson:"SystemSpecification"` // 1 for POD, 2 for ECS
	ImageUrl            string   `json:"ImageUrl,omitempty" bson:"ImageUrl"`                       // Docker image URL for POD, system image for ECS
	ReScript            string   `json:"ReScript,omitempty" bson:"ReScript"`                       // Used to estimate the resources consumed by the task
	ScScript            string   `json:"ScScript,omitempty" bson:"ScScript"`                       // Used to collect task status and progress from the output of the command line
	JeScript            string   `json:"JeScript,omitempty" bson:"JeScript"`                       // Task entry command
	FixedParameters     []string `json:"FixedParameters,omitempty" bson:"FixedParameters"`         // Relevant parameters for the task
}

// Job Type Filter for Different Database Queries
type JobTypeFilter struct {
	JobTypeIdFilter           string   `json:"JobTypeIdFilter"`
	SystemSpecificationFilter int32    `json:"SystemSpecificationFilter"`
	ImageUrlFilter            string   `json:"ImageUrlFilter"`
	FixedParametersFilter     []string `json:"FixedParametersFilter"`
	ReScriptFilter            string   `json:"ReScriptFilter"`
	ScScriptFilter            string   `json:"ScScriptFilter"`
	JeScriptFilter            string   `json:"JeScriptFilter"`
}
