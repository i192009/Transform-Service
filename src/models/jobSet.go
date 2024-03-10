package models

// JobSet represents a set of jobs with associated metadata.
type JobSet struct {
	JobSetId        string `json:"JobSetId,omitempty" bson:"JobSetId,omitempty"`     // Unique identifier for the job set
	JobTypeId       string `json:"JobTypeId,omitempty" bson:"JobTypeId"`             // Identifier for the type of jobs in the set
	Name            string `json:"Name,omitempty" bson:"Name"`                       // Name of the job set in the database
	Total           int64  `json:"Total,omitempty" bson:"Total"`                     // Total number of jobs in the set
	FixedParameters string `json:"FixedParameters,omitempty" bson:"FixedParameters"` // Fixed parameters associated with the job set
}

// JobSetFilter represents filters for querying JobSets.
type JobSetFilter struct {
	JobSetIdFilter        string `json:"JobSetIdFilter"`        // Filter by resource pool ID
	JobTypeIdFilter       string `json:"JobTypeIdFilter"`       //Filter by Namespace
	NameFilter            string `json:"NameFilter"`            // Filter by resource pool name
	TotalFilter           string `json:"TotalFilter"`           // Filter by whether the resource pool is shared
	FixedParametersFilter string `json:"FixedParametersFilter"` // Filter by scaling strategy
}
