package service

import (
	"github.com/hdzattain/smart-gateway/setting/operation_setting"
	"github.com/hdzattain/smart-gateway/setting/system_setting"
)

func GetCallbackAddress() string {
	if operation_setting.CustomCallbackAddress == "" {
		return system_setting.ServerAddress
	}
	return operation_setting.CustomCallbackAddress
}
