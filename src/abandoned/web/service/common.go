package service

import (
	"fmt"
	"time"
	"transform2/config"

	"gitlab.zixel.cn/go/framework/bus"
	"gitlab.zixel.cn/go/framework/logger"
)

type BusMail struct {
	Subject  string
	Body     string
	Nickname string
	Tos      []string
}

var log = logger.Get()

const (
	mtRecordRedisKeyPrefix = "S:filetransfer:%v:%v:%v"
	recordExpire           = 10 * 60 * time.Second
)
const (
	lockRedisKeyPrefix = "S:filetransfer:lock:%v:C"
	delCommand         = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
)

func _InitRedisKey(keyPrefix, key, serviceMark, paramType string) string {
	return fmt.Sprintf(keyPrefix, serviceMark, key, paramType)
}

func Key(key, serviceMark, paramType string) string {
	return _InitRedisKey(mtRecordRedisKeyPrefix, key, serviceMark, paramType)
}

func _InitRedisLockKey(keyPrefix, key string) string {
	return fmt.Sprintf(keyPrefix, key)
}

func LockKey(key string) string {
	return _InitRedisLockKey(lockRedisKeyPrefix, key)
}

var ConfigEventBus = bus.NewAsyncEventBus()

const (
	Topic_Mail = "Topic_Mail"
)

func init() {
	ConfigEventBus.Subscribe(Topic_Mail, func(data BusMail) {
		mailAddress := config.MailTo["Job"].([]string)
		if len(data.Tos) > 0 {
			mailAddress = data.Tos
		}
		mail := config.Mail{
			MailHost: config.MailHost,
			MailPort: config.MailPort,
			MailUser: config.MailUser,
			MailPwd:  config.MailPassword,
		}
		if len(data.Nickname) <= 0 {
			data.Nickname = "<" + config.MailUser + ">"
		}
		Subject := "env:" + config.Namespace + ">>" + data.Subject
		if err := config.SendGoMail(mailAddress, Subject, data.Body, data.Nickname, mail); err != nil {
			log.Errorf("ConfigEventBus is  error,topic:%v,data:%v,err:%v", data.Subject, data.Body, err)
		}
	})
}
