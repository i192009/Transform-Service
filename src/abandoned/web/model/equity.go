package model

import "time"

type Equity struct {
	EquityId     string    `bson:"_id" validate:"required" json:"equityId"`
	UserId       string    `bson:"userId" json:"userId"`
	TenantId     string    `bson:"tenantId" json:"tenantId"`
	Weight       int       `bson:"weight" json:"weight"`
	UserType     int       `bson:"userType" json:"userType"`
	Marks        []string  `bson:"marks" json:"marks"`
	AvailableNum int       `bson:"availableNum" json:"availableNum"`
	StarTime     time.Time `bson:"starTime" json:"starTime"`
	EndTime      time.Time `bson:"endTime" json:"endTime"`
	IsEnable     bool      `bson:"isEnable" json:"isEnable"`
	UpdateTime   time.Time `bson:"updateTime" json:"updateTime"`
	createTime   time.Time `bson:"createTime" json:"createTime"`
}
