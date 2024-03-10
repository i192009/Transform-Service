package model

type NewJobReq struct {
	//ObjectKey     string              `json:"objectKey" validate:"required"`   /// directoryOrFileToBeConverted
	Bucket_DL    string          `json:"bucketDownload,omitempty"`         /// The bucket where the file to be converted is located
	Bucket_UL    string          `json:"bucketUpload,omitempty"`           /// bucketUploadedAfterFileConversionIsCompleted
	FolderConfig []*FileConfig   `json:"folderConfig" validate:"required"` /// File List
	Scale        float32         `json:"scale,omitempty"`                  /// The conversion coefficient from user units to millimeters. For example, if the user's model unit is "meter", fill in 1000 here.
	TargetFormat []*TargetFormat `json:"targetFormat" validate:"required"` /// conversionTargetFormat
	Pipelines    []*Pipeline     `json:"pipelines,omitempty"`
	ServiceType  int             `json:"serviceType,omitempty"` /// _1_hs_2_pz_is_used_by_default
	Marks        []string        `json:"marks"`                 /// Mark data environment dev development environment release is the formal environment
	JobMqInfo    *JobMqInfo      `json:"jobMqInfo"`
}
