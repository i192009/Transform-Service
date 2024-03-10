package model

var (
	PROCESS_PENDING          = 0
	PROCESS_START            = 1
	PROCESS_DOWNLOAD         = 2
	PROCESS_DOWNLOAD_FAILED  = 3
	PROCESS_PROCESSING       = 4
	PROCESS_CANCELING        = 5
	PROCESS_CANCELED         = 6
	PROCESS_PROCESSED        = 7
	PROCESS_PROCESS_FAILED   = 8
	PROCESS_UPLOADING        = 9
	PROCESS_UPLOADING_FAILED = 10
	PROCESS_TIMEOUT          = 11 //调度服务自定义超时状态，不更新到数据库
)

var EndStatusMap = map[int]interface{}{
	PROCESS_PROCESSED:        nil,
	PROCESS_PROCESS_FAILED:   nil,
	PROCESS_UPLOADING_FAILED: nil,
	PROCESS_TIMEOUT:          nil,
	PROCESS_CANCELED:         nil,
	PROCESS_DOWNLOAD_FAILED:  nil,
}

var FailedStatusMap = map[int]interface{}{
	PROCESS_DOWNLOAD_FAILED:  nil,
	PROCESS_PROCESS_FAILED:   nil,
	PROCESS_UPLOADING_FAILED: nil,
}
