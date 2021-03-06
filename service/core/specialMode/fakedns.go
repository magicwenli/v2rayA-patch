package specialMode

import (
	"github.com/magicwenli/v2rayA-patch/common"
	"github.com/magicwenli/v2rayA-patch/core/v2ray/where"
	"github.com/magicwenli/v2rayA-patch/db/configure"
)

func CouldUseFakeDns() bool {
	variant, ver, err := where.GetV2rayServiceVersion()
	if err != nil {
		ver = "0.0.0"
	}
	if variant == where.V2ray {
		if fakeDnsValid, _ := common.VersionGreaterEqual(ver, "4.35.0"); !fakeDnsValid {
			return false
		}
	}
	return configure.GetSettingNotNil().AntiPollution != configure.AntipollutionClosed
}

func ShouldUseFakeDns() bool {
	return configure.GetSettingNotNil().SpecialMode == configure.SpecialModeFakeDns
}
