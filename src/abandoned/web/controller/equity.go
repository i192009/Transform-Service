package controller

import (
	"transform2/abandoned/web/model"
	"transform2/abandoned/web/service"
	"transform2/config"
)

func CheckEquity(userID, tenantID *string, runNum int) (bool, error) {
	//添加分布式锁

	//添加分布式锁
	EquityDBList, err := service.GetEquityListByUserIDOrTenantId(userID, tenantID)
	if err != nil {
		return false, err
	}
	if len(EquityDBList) <= 0 {
		return true, nil
	}
	for _, equityDB := range EquityDBList {
		if !equityDB.IsEnable {
			continue
		}
		if equityDB.AvailableNum >= runNum {
			equityDB.AvailableNum = equityDB.AvailableNum - runNum
			service.UpdateEquity(equityDB.EquityId, equityDB)
			return true, nil
		} else {
			return false, nil
		}
	}
	return true, nil
}

func GetTaskScopeId(userinfo *model.UserInfo) (*int32, string) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("GetTaskScopeId>>>panic:%v", err)
		}
	}()
	if userinfo == nil {
		log.Errorf("GetTaskScopeId>>>userinfo is nil")
		return nil, ""
	}
	equityConfigMap := config.EquityConfigMap
	log.Infof("GetTaskScopeId>>>equityConfigMap:%v", equityConfigMap)
	if equityConfigMap == nil || len(equityConfigMap) <= 0 {
		log.Errorf("GetTaskScopeId>>>equityConfigMap is nil")
		return nil, ""
	}
	if userinfo.TenantId != nil {
		if equityConfig, ok := equityConfigMap[*userinfo.TenantId].(map[string]int32); ok {
			scopeId := equityConfig["scope_id"]
			return &scopeId, *userinfo.TenantId
		}
	}
	if userinfo.UserId != nil {
		if equityConfig, ok := equityConfigMap[*userinfo.UserId].(map[string]int32); ok {
			scopeId := equityConfig["scope_id"]
			return &scopeId, *userinfo.UserId
		}
	}
	return nil, ""
}

func CheckTaskTargetFormat(userinfo *model.UserInfo, targetFormats []string) bool {
	if userinfo == nil {
		return true
	}
	equityConfigMap := config.EquityConfigMap
	if equityConfigMap == nil || len(equityConfigMap) <= 0 {
		return true
	}
	var supportFormat []string
	if userinfo.TenantId != nil {
		if equityConfig, ok := equityConfigMap[*userinfo.TenantId].(map[string]interface{}); ok {
			if len(equityConfig["support_format"].([]string)) == 0 {
				return true
			}
			supportFormat = equityConfig["support_format"].([]string)
		}
	}
	if userinfo.UserId != nil && len(supportFormat) == 0 {
		if equityConfig, ok := equityConfigMap[*userinfo.UserId].(map[string]interface{}); ok {
			if len(equityConfig["support_format"].([]string)) == 0 {
				return true
			}

			supportFormat = equityConfig["support_format"].([]string)
		}
	}
	if len(supportFormat) > 0 {
		supportFormatMap := make(map[string]bool)
		for _, format := range supportFormat {
			supportFormatMap[format] = true
		}
		for _, format := range targetFormats {
			if _, ok := supportFormatMap[format]; !ok {
				return false
			}
		}
	}
	return true
}

func CheckVIP(userID, tenantID *string) bool {
	Equity_DB, err := service.GetEquityListByUserIDOrTenantId(userID, tenantID)
	if err != nil {
		return false
	}
	if len(Equity_DB) <= 0 {
		return false
	}
	for _, equityDB := range Equity_DB {
		if !equityDB.IsEnable {
			continue
		}
		if equityDB.UserType == 2 {
			return true
		}
	}
	return false
}
