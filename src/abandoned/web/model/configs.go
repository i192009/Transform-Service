package model

type FileConfig struct {
	Name         string                 `json:"name,omitempty"`
	Transform    bool                   `json:"transform,omitempty"` //Whether file conversion is required
	FileSize     int64                  `json:"fileSize,omitempty"`  //File size in m
	RemoteUrl    string                 `json:"remoteUrl,omitempty"`
	TargetUrl    map[string]string      `json:"targetUrl,omitempty"` //Address after conversion
	TargetUrlMap map[string]string      `json:"targetUrlMap"`        //Target format, corresponding to the target file address, does not convert and does not fill in the content, but the fields must be retained
	Pip          map[string]interface{} `json:"p,omitempty"`
}

type Pipeline struct {
	UseScripts  int             `json:"useScripts"` ////1: Use script processing.
	ProcessType int             `json:"processType,omitempty"`
	Parms       *PipelineParams `json:"parms"`
}

type PipelineParams struct {
	Scripts     []string      `json:"scripts,omitempty"`
	Optimize    *OptimizeFunc `json:"optimize,omitempty"`    /// Optimization parameters
	ProcessMesh *ProcessMesh  `json:"processMesh,omitempty"` /// Surface reduction parameters
}

type TargetFormat struct {
	Name  string      `json:"name,omitempty"`
	Pipe  []int       `json:"pipe,omitempty"`
	Parms interface{} `json:"parms,omitempty"`
	Tag   string      `json:"tag,omitempty"`
}
