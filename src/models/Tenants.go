package models

// Config Tenant struct to init the configs and pools
type TenantConfig struct {
	ConfigId        string       `json:"ConfigId"`
	Set             JobSet       `json:"Set"`
	Pool            ResourcePool `json:"Pool"`
	AssociatedPools []string     `json:"AssociatedPools"`
}
