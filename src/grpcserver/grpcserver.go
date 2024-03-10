package grpcserver

import (
	"transform2/config"
	"transform2/services"

	"gitlab.zixel.cn/go/framework"
	"gitlab.zixel.cn/go/framework/logger"
	"go.temporal.io/sdk/client"
)

type TransformServer struct {
	services.UnimplementedTransformV2Server
	WorkflowClient client.Client
}

var log = logger.Get()

func Init() error {
	c, err := client.Dial(client.Options{
		Namespace: config.TemporalNamespace,
		HostPort:  config.TemporalAddress,
	})

	if err != nil {
		log.Fatalln("Unable to create client", err)
		return err
	}

	framework.RegisterService(&services.TransformV2_ServiceDesc, &TransformServer{WorkflowClient: c})
	framework.RegisterService(&services.JobManagement_ServiceDesc, &JobManageServer{})
	framework.RegisterService(&services.ResourcePoolManagement_ServiceDesc, &ResourcePoolServer{})
	framework.RegisterService(&services.TenantManagement_ServiceDesc, &TenantConfigServer{})
	config.InitMongoDB()
	return nil
}
