package models

// ResourcePool represents a pool of resources with associated configuration.
type ResourcePool struct {
	ResourcePoolID  string                          `json:"ResourcePoolID" bson:"ResourcePoolId"`             // Unique identifier for the resource pool
	NameSpace       string                          `json:"NameSpace" bson:"NameSpace"`                       //Temporal Namespace
	Name            string                          `json:"Name,omitempty" bson:"Name"`                       // Name of the resource pool; uses `omitempty` to omit empty values in JSON
	IsShared        bool                            `json:"IsShared,omitempty" bson:"IsShared"`               // Indicates whether the resource pool is shared
	ScalingStrategy int32                           `json:"ScalingStrategy,omitempty" bson:"ScalingStrategy"` // Scaling strategy for the resource pool
	ScalingLimit    int32                           `json:"ScalingLimit,omitempty" bson:"ScalingLimit"`       // Scaling limit for the resource pool
	QueueLimit      int32                           `json:"QueueLimit,omitempty" bson:"QueueLimit"`           // Queue limit for the resource pool
	Fixed           int32                           `json:"Fixed,omitempty" bson:"Fixed"`                     // Fixed value associated with the resource pool
	DefaultTaskSet  int32                           `json:"DefaultTaskSet,omitempty" bson:"DefaultTaskSet"`   // Default task set for the resource pool
	ResourceLimits  map[string]*ResourceLimitOfTask `json:"ResourceLimits,omitempty" bson:"ResourceLimits"`   // List of resource limits associated with tasks
}

// ResourceLimitOfTask represents resource limits for a specific task type.
type ResourceLimitOfTask struct {
	JobTypeId    string `json:"JobTypeId,omitempty"`           // Identifier for the type of task
	ScalingLimit int32  `json:"PercentScalingLimit,omitempty"` //percentage scaling limit
}

// ResourcePoolFilter represents filters for querying resource pools.
type ResourcePoolFilter struct {
	ResourcePoolIDFilter  string `json:"resourcePoolIDFilter"`  // Filter by resource pool ID
	NameSpaceFilter       string `json:"nameSpaceFilter"`       //Filter by Namespace
	NameFilter            string `json:"nameFilter"`            // Filter by resource pool name
	IsSharedFilter        bool   `json:"isSharedFilter"`        // Filter by whether the resource pool is shared
	ScalingStrategyFilter int32  `json:"scalingStrategyFilter"` // Filter by scaling strategy
	ScalingLimitFilter    int32  `json:"scalingLimitFilter"`    // Filter by scaling limit
	CustomFilter          string `json:"customFilter"`          // Custom filter for additional criteria
}
