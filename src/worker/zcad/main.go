package main

import (
	"transform2/worker/zcad/zcadworker"

	"gitlab.zixel.cn/go/framework/logger"
)

var log = logger.Get()

func main() {
	zcadworker.Start()
}
