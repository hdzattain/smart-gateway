package controller

import (
	"strings"

	"github.com/hdzattain/smart-gateway/common"
	"github.com/hdzattain/smart-gateway/setting/system_setting"
)

func paymentReturnPath(suffix string) string {
	base := strings.TrimRight(system_setting.ServerAddress, "/")
	return base + common.ThemeAwarePath(suffix)
}
