package constants

import "github.com/spf13/viper"

var (
	HomeURLBase   = viper.GetString("homed.url")
	DevicesEp     = HomeURLBase + "/devices"
	DevicesPingEp = DevicesEp + "/ping"
	AssetsEp      = HomeURLBase + "/assets"
	GetAssetEp    = AssetsEp + "/%d/download"
)
