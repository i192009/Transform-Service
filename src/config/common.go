package config

import (
	"time"

	"github.com/joho/godotenv"
	"gitlab.zixel.cn/go/framework/config"
)

var (
	Timeout     = 600 * time.Second
	_           = godotenv.Load()
	Namespace   = config.GetString("namespace", "default")
	ServiceName = "transform"
	ServiceUUID = "UUID-af7455-98ab372-652c1bb"

	ListenAddr = config.GetString("listen_addr", "")
	TaskLimit  = config.GetInt("task_limit", 0)
	JobTimeout = int(config.GetInt("job_timeout", 0))

	MailHost     = config.GetString("mail.host", "")
	MailPort     = int(config.GetInt("mail.port", 0))
	MailUser     = config.GetString("mail.user", "")
	MailPassword = config.GetString("mail.password", "")
	MailTo       = config.GetObject("mail.to")

	HuaweiEnable      = config.GetBoolean("huawei.Enable", false)
	HuaweiAccessKey   = config.GetString("huawei.AK", "")
	HuaweiSecretKey   = config.GetString("huawei.SK", "")
	HuaweiRegion      = config.GetString("huawei.Region", "")
	HuaweiImageRef    = config.GetString("huawei.ImageRef", "")
	HuaweiSubnetId    = config.GetString("huawei.SubnetId", "")
	HuaweiFlavorRef   = config.GetString("huawei.FlavorRef", "")
	HuaweiVpcId       = config.GetString("huawei.VpcId", "")
	HuaweiMaxNum      = config.GetInt("huawei.MaxNum", 0)
	HuaweiPollingTime = config.GetInt("huawei.PollingTime", 0)
	HuaweiWaitTime    = float64(config.GetInt("huawei.WaitTime", 0))
	HuaweiWaitNum     = int(config.GetInt("huawei.WaitNum", 0))
	HuaweiExpansion   = int(config.GetInt("huawei.Expansion", 0))

	TopicOperationLogMap = config.GetObject("mq_config.topic_config.topicOperationLog.service_key_map")
	TopicJobProgressMap  = config.GetObject("mq_config.topic_config.topicJobProgress.service_key_map")

	ClearJobConfigAppIds      = config.GetArray("clear_job_config.appId")
	ClearJobConfigInstanceIds = config.GetArray("clear_job_config.instanceId")
	ClearJobConfigTenantIds   = config.GetArray("clear_job_config.tenantId")
	ClearJobConfigUserIds     = config.GetArray("clear_job_config.userId")

	TemporalAddress   = config.GetString("temporal.address", "localhost:7233")
	TemporalNamespace = config.GetString("temporal.namespace", "default")

	EquityConfigMap = config.GetObject("equity_config")

	Scripts = []string{
		"",
		"algo.retessellate([1], 5, -1, -1);scene.mergeFinalLevel([1], scene.MergeHiddenPartsMode.MergeSeparately, True);material.makeMaterialNamesUnique();algo.removeHoles([1], True, True, True, 50.0);algo.deletePatches([]);algo.decimate([1], 5.0, 0.1, 1.0, -1.0, False)",                                                                                                                 //Kreat Component
		"scene.mergePartsByAssemblies([1],2);scene.selectByMaximumSize([1],150);scene.deleteSelection();algo.retessellate([1], 10.000000, -1, -1);scene.mergeFinalLevel([1], 2, True);material.makeMaterialNamesUnique();algo.removeHoles([1], True, True, True, 50.0);algo.deletePatches([]);algo.decimate([1], 10.0, 0.1, 1.0, -1.0, False);scene.deleteEmptyOccurrences();scene.compress(1)", //Kreat full
	}
)
