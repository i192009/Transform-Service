package controller

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"transform2/abandoned/web/model"
	"transform2/config"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.zixel.cn/go/framework/database"
	"gitlab.zixel.cn/go/framework/logger"
)

var log = logger.Get()
var validate = validator.New()

func NewTask(ctx context.Context, userinfo *model.UserInfo, req *model.NewJobReq, scripts []string) (string, *config.ErrorNo) {
	log.Infof("NewTask req:%s", req.Pipelines)
	if isPass, err := CheckEquity(userinfo.UserId, userinfo.TenantId, 1); !isPass && err == nil {
		return "", config.NewErrorNo(30101, "CheckEquity is error", nil)
	}

	// sID, _ := strconv.Atoi(*userinfo.ScopeId)

	// conn := framework.GetGrpcConnection("StorageService")
	// grpc := services.NewStorageServiceClient(conn)

	// dl_buck, _ := grpc.GetDownLoadSignUrl(ctx, &services.DownLoadSignUrlRequest{
	// 	ScopeId: int32(sID),
	// 	Key:     config.ObsDlBucket,
	// })
	// ul_buck, _ := grpc.GetUploadFileUrl(ctx, &services.UploadFileUrlRequest{
	// 	ScopeId: int32(sID),
	// 	Key:     config.ObsUlBucket,
	// })

	// req.Bucket_UL = ul_buck.Url
	// req.Bucket_DL = dl_buck.Url

	//判断目标格式是否可用
	targetFormats := make([]string, len(req.TargetFormat))
	for index, format := range req.TargetFormat {
		targetFormats[index] = format.Name
	}
	if !CheckTaskTargetFormat(userinfo, targetFormats) {
		log.Error("CheckTaskTargetFormat error")
		return "", config.NewErrorNo(30102, "", nil)
	}
	//分布式锁
	if err := validate.Struct(req); err != nil {
		log.Error(err.Error())
		return "", config.NewErrorNo(30001, "Invalid request format.", err)
	}
	data, _ := json.Marshal(req)
	lockKey := config.MD5(string(data))
	randomStr := config.RandomString(16)
	timeout := time.Duration(10) * time.Second
	if err := database.RedisSetCtx(ctx, lockKey, randomStr, &timeout); err != nil {
		log.Error("newjobV2 分布式锁-请勿重复请求")
		return "", config.NewErrorNo(30004, "请勿重复请求", nil)
	}
	defer func(ctx context.Context, keys []string) {
		err := database.RedisMultiDelCtx(ctx, keys)
		if err != nil {
			log.Error("newjobV2 分布式锁-删除失败")
		}
	}(ctx, []string{lockKey})
	if req.FolderConfig == nil || len(req.FolderConfig) <= 0 {
		log.Error("FolderConfig is nil")
		return "", config.NewErrorNo(30002, "", nil)
	}
	jobId := uuid.NewString()
	initJob := func() *model.Job {
		job := new(model.Job)
		job.JobId = jobId
		job.JobStatus.Status = model.PROCESS_PENDING
		job.JobStatus.Total = 0
		job.JobStatus.Processed = 0
		job.JobStatus.Progress = 0
		job.ClientId = 0xfffe
		return job
	}
	job := initJob()
	job.AppId = userinfo.AppId
	job.UserId = userinfo.UserId
	job.InstanceId = userinfo.InstanceId
	job.TenantId = userinfo.TenantId
	if job == nil {
		return "", config.NewErrorNo(30100, "", nil)
	}
	log.Info("newjobV3 init,job:%s", jobId)
	job.UploadUrl = filepath.Join("converted", job.RemoteUrl)
	job.Version = "3.0.0"

	/// 设置上传和下载的桶
	if len(req.Bucket_DL) > 0 {
		job.BucketDl = req.Bucket_DL
	}

	if len(req.Bucket_UL) > 0 {
		job.BucketUl = req.Bucket_UL
	}
	for _, pipeline := range req.Pipelines {
		pipeline.Parms.Scripts = scripts
	}
	job.Pipeline = req.Pipelines
	job.TargetFormat = req.TargetFormat
	job.FolderConfig = req.FolderConfig
	job.ServiceType = req.ServiceType
	job.Marks = req.Marks
	job.JobMqInfo = req.JobMqInfo
	/// 转为目标格式
	DownloadMaps := make(map[string]map[string]string)
	for _, fileConfig := range job.FolderConfig {
		if !fileConfig.Transform {
			continue
		}
		if len(fileConfig.RemoteUrl) == 0 {
			continue
		}

		//fileMap:=make(map[string]string)
		uploadMap := make(map[string]string)
		for index, targetFormat := range job.TargetFormat {
			var (
				uploadUrl    string
				targetUrlTag = fileConfig.TargetUrl[targetFormat.Tag]
			)
			if targetUrl, ok := fileConfig.TargetUrlMap[strconv.FormatInt(int64(index), 10)]; ok {
				uploadUrl = filepath.Join(
					targetUrlTag, targetUrl)
			} else {
				uploadUrl = filepath.ToSlash(
					filepath.Join(
						targetUrlTag,
						//job.RemoteUrl,
						strings.TrimSuffix(filepath.Base(fileConfig.Name), filepath.Ext(fileConfig.Name))+"."+targetFormat.Name,
					),
				)
			}
			job.Uploads = append(job.Uploads, uploadUrl)
			uploadMap[targetFormat.Tag] = uploadUrl
		}
		DownloadMaps[fileConfig.RemoteUrl] = uploadMap
		job.UploadMaps = append(job.UploadMaps, uploadMap)
		job.DownloadMaps = DownloadMaps
	}
	log.Infof("newjobV3_1 PostEvent,marks:%v,job:%s", job.Marks, jobId)
	nowTime := time.Now()
	job.CreateTime = nowTime
	if CheckVIP(userinfo.UserId, userinfo.TenantId) {
		job.JobType = 2
		model.ServiceServer.PostEventVIP(model.EVENT_CREATE_JOB, job)
	} else {
		model.ServiceServer.PostEvent(model.EVENT_CREATE_JOB, job)
	}
	return jobId, nil
}
