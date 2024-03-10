package model

import "time"

type Job struct {
	UserId        *string
	InstanceId    *string
	AppId         *string
	TenantId      *string
	Version       string                       /// requestedVersionNumber
	JobId         string                       /// jobID
	JobStatus     JobStatus                    /// taskExecutionStatus
	ClientId      uint16                       /// clientId
	BucketDl      string                       /// theBucketWhereTheFileIsLocated
	BucketUl      string                       /// bucketUsedForUploading
	RemoteUrl     string                       /// downloadPath
	UploadUrl     string                       /// uploadPath
	Tasks         []string                     /// pendingTasks
	Uploads       []string                     /// uploadPath
	UploadMaps    []map[string]string          /// uploadedPathMap   key: TargetFormat  value:UploadUrl
	DownloadMaps  map[string]map[string]string /// downloadPathMap   key: TargetFormat  value:DownloadUrl
	TargetFormats []string                     /// targetFormat
	ServiceType   int                          /// 1,HS，usedByDefault2,PZ
	Scale         float32                      /// The conversion coefficient from user units to millimeters. For example, if the user's model unit is "meter", fill in 1000 here.
	JobType       int                          //1,usually  2，VIP
	OptimizeFuns  OptimizeFunc                 /// optimizationParameters
	ProcessMesh   ProcessMesh                  /// surfaceReductionParameters

	RetryCount   int           /// numberOfRetries
	Marks        []string      /// markupDataEnvironment  dev developmentEnvironment  release isAFormalEnvironment
	FolderConfig []*FileConfig /// uploadFolderContents
	Scripts      []string      //scenario
	UseScript    int           ////1:useScriptProcessing

	Pipeline     []*Pipeline
	TargetFormat []*TargetFormat
	JobMqInfo    *JobMqInfo
	CreateTime   time.Time
}

type JobStatus struct {
	Status     int                    /// taskStatus
	Total      int                    /// totalNumberOfTasks
	Processed  int                    /// numberOfTasksProcessed
	Progress   int                    /// processingProgress
	Time       float64                //processingTime
	Result     map[string]interface{} /// numberOfSidesProcessed
	UpdateTime time.Time              /// lastUpdatedTime
}

type JobMqInfo struct {
	Tag          string `json:"mq_tag"`
	SendRateType string `json:"send_rate_type"` //1, progress  2. Status
}
