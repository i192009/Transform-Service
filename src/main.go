package main

import (
	"transform2/grpcserver"
	"transform2/web"

	"gitlab.zixel.cn/go/framework"
)

func main() {
	grpcserver.Init()
	framework.LoadServiceRoute(web.SetupProbeRoutes, "v2")
	framework.Run()
}
