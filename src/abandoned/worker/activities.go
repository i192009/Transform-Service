package main

// func NewTask(ctx context.Context, in *services.C2S_NewTaskReqT) (*services.C2S_NewTaskRpnT, error) {

// 	controller.LimitJobHandler() //Call the Job Handler Function before processing the incoming Request

// 	log.Info("Activity Started")
// 	defer func() {
// 		if err := recover(); err != nil {
// 			log.Errorf("NewTask panic:%v", err)
// 		}
// 	}()
// 	by, _ := json.Marshal(in)
// 	log.Infof("C2S_NewTaskReqT.P:%s", string(by))
// 	folderConfig := make([]*model.FileConfig, len(in.FolderConfig))
// 	for index, cfg := range in.FolderConfig {
// 		pip := make(map[string]interface{})
// 		for key, value := range cfg.P {
// 			pip[key] = config.StructpbToMap(value)
// 		}

// 		folderConfig[index] = &model.FileConfig{
// 			Name:         cfg.Name,
// 			Transform:    cfg.Transform,
// 			FileSize:     cfg.FileSize,
// 			RemoteUrl:    cfg.RemoteUrl,
// 			TargetUrl:    cfg.TargetUrl,
// 			TargetUrlMap: cfg.TargetUrl,
// 			Pip:          pip,
// 		}
// 	}
// 	targetFormat := make([]*model.TargetFormat, len(in.TargetFormats))
// 	scripts := make([]string, len(in.TargetFormats))
// 	for index, format := range in.TargetFormats {
// 		pipe := make([]int, len(format.Pipe))
// 		for i, p := range format.Pipe {
// 			pipe[i] = int(p)
// 		}
// 		parms := config.StructpbToMap(format.Parms)
// 		log.Debugf("parms:%v", parms)
// 		_, script := controller.CheckTargetFormatPram(parms)
// 		log.Debugf("script:%s", script)
// 		//{"0":  {  "useScript":2 } }
// 		targetFormatParam := map[string]map[string]int{"0": {"useScript": index}}
// 		//script, targetFormatParam := logic.GetScripts(parms)
// 		scripts[index] = script
// 		targetFormat[index] = &model.TargetFormat{
// 			Name:  format.Name,
// 			Pipe:  pipe,
// 			Parms: targetFormatParam,
// 			Tag:   format.Tag,
// 		}
// 	}

// 	req := &model.NewJobReq{
// 		Bucket_DL:    in.BucketDownload,
// 		Bucket_UL:    in.BucketUpload,
// 		FolderConfig: folderConfig,
// 		//Scale: in.Scale,
// 		TargetFormat: targetFormat,
// 		//Pipelines
// 		ServiceType: int(in.ServiceType),
// 		Marks:       in.Marks,
// 	}
// 	pipelines := make([]*model.Pipeline, len(in.Pipelines))
// 	for index, pipeline := range in.Pipelines {
// 		pipelines[index] = &model.Pipeline{
// 			UseScripts:  int(pipeline.UseScripts),
// 			ProcessType: int(pipeline.ProcessType),
// 			Parms:       &model.PipelineParams{},
// 		}
// 	}
// 	req.Pipelines = pipelines
// 	if in.JobMqInfos != nil {
// 		req.JobMqInfo = &model.JobMqInfo{
// 			Tag:          in.JobMqInfos.MqTag,
// 			SendRateType: in.JobMqInfos.SendRateType,
// 		}
// 	}
// 	userinfo := &model.UserInfo{
// 		UserId:     in.UserId,
// 		InstanceId: in.InstanceId,
// 		AppId:      in.AppId,
// 		TenantId:   in.TenantId,
// 		ScopeId:    in.ScopeID,
// 	}
// 	if jobId, errn := controller.NewTask(ctx, userinfo, req, scripts); errn != nil {
// 		return &services.C2S_NewTaskRpnT{
// 			Code: int32(errn.Code),
// 			Msg:  errn.ErrMsg(),
// 		}, nil
// 	} else {
// 		return &services.C2S_NewTaskRpnT{
// 			Code:  0,
// 			Msg:   "success",
// 			JobId: jobId,
// 		}, nil
// 	}
// }
