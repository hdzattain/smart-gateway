package system_setting

import (
	"github.com/hdzattain/smart-gateway/common"
	"github.com/hdzattain/smart-gateway/setting/config"
)

type ThemeSettings struct {
	Frontend string `json:"frontend"`
}

var themeSettings = ThemeSettings{
	Frontend: "classic",
}

func init() {
	config.GlobalConfig.Register("theme", &themeSettings)
	syncThemeToCommon()
}

func syncThemeToCommon() {
	common.SetTheme(themeSettings.Frontend)
}

func GetThemeSettings() *ThemeSettings {
	return &themeSettings
}

// UpdateAndSyncTheme syncs the theme config to common after DB load.
func UpdateAndSyncTheme() {
	syncThemeToCommon()
}
