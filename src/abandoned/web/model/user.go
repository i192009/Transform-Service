package model

type UserInfo struct {
	UserId     *string `json:"userId,omitempty"`
	InstanceId *string `json:"instanceId,omitempty"`
	AppId      *string `json:"appId,omitempty"`
	TenantId   *string `json:"tenantId,omitempty"`
	ScopeId    *string `json:"scopeId,omitempty"`
}
