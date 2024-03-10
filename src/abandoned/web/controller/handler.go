package controller

import (
	"fmt"
	"time"
	"transform2/abandoned/web/model"
	"transform2/abandoned/web/service"

	"github.com/gin-gonic/gin"
	"gitlab.zixel.cn/go/framework"
	"gitlab.zixel.cn/go/framework/database"
	"gitlab.zixel.cn/go/framework/logger"
)

func LimitJobHandler() gin.HandlerFunc {
	var log = logger.Get()
	return func(context *gin.Context) {
		if len(model.ServiceServer.JobQueue) >= 1000 {
			log.Errorf("too many tasks, try again later.")
			framework.Error(context, framework.NewServiceError(30201, "too many tasks, try again later."))
			context.Abort()
		}
		//When the number of queued tasks exceeds 300, a letter needs to be sent to the operation and maintenance team,
		//and the operation and maintenance team will decide whether to increase the background computing power.
		//The interval between two emails should be no less than 30 minutes
		if len(model.ServiceServer.JobQueue) > 500 {
			redisKey := service.Key("limit", "job.data", "C")
			if ok, _ := database.RedisExists(context.Request.Context(), redisKey); !ok {
				body := "<html>\n\t\t\t\t\t\t\t<body>\n\t\t\t\t\t\t\t\t<h1>" +
					"当前排队任务： " + fmt.Sprintf("%d", len(model.ServiceServer.JobQueue)) + "\n" +
					"请注意！！！" +
					"</h1>\n\t\t\t\t\t\t\t</body>\n\t\t\t\t\t\t</html>"
				data := service.BusMail{
					Subject: "任务预警",
					Body:    body,
				}
				service.ConfigEventBus.Publish(service.Topic_Mail, data)
				timeout := time.Duration(30) * time.Minute
				err := database.RedisSetCtx(context.Request.Context(), redisKey, "", &timeout)
				if err != nil {
					log.Errorf(err.Error())
					return
				}
			}
		}
		//context.Next()
	}
}
