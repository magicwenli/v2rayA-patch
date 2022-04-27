package v2ray

import (
	"fmt"
	"strings"
	"time"

	"github.com/magicwenli/v2rayA-patch/conf"
	"github.com/magicwenli/v2rayA-patch/core/iptables"
	"github.com/magicwenli/v2rayA-patch/core/specialMode"
	"github.com/magicwenli/v2rayA-patch/db/configure"
	"github.com/magicwenli/v2rayA-patch/pkg/util/log"
)

func DeleteTransparentProxyRules() {
	iptables.CloseWatcher()
	if !conf.GetEnvironmentConfig().Lite {
		removeResolvHijacker()
		iptables.Tproxy.GetCleanCommands().Run(false)
		iptables.Redirect.GetCleanCommands().Run(false)
		iptables.DropSpoofing.GetCleanCommands().Run(false)
	}
	iptables.SystemProxy.GetCleanCommands().Run(false)
	time.Sleep(100 * time.Millisecond)
}

func WriteTransparentProxyRules() (err error) {
	defer func() {
		if err != nil {
			log.Warn("WriteTransparentProxyRules: %v", err)
			DeleteTransparentProxyRules()
		}
	}()
	if specialMode.ShouldUseSupervisor() {
		if err = iptables.DropSpoofing.GetSetupCommands().Run(true); err != nil {
			log.Warn("DropSpoofing can't be enable: %v", err)
			return err
		}
	}
	setting := configure.GetSettingNotNil()
	switch setting.TransparentType {
	case configure.TransparentTproxy:
		if err = iptables.Tproxy.GetSetupCommands().Run(true); err != nil {
			if strings.Contains(err.Error(), "TPROXY") && strings.Contains(err.Error(), "No chain") {
				err = fmt.Errorf("you does not compile xt_TPROXY in kernel")
			}
			return fmt.Errorf("not support \"tproxy\" mode of transparent proxy: %w", err)
		}
		iptables.SetWatcher(&iptables.Tproxy)
	case configure.TransparentRedirect:
		if err = iptables.Redirect.GetSetupCommands().Run(true); err != nil {
			return fmt.Errorf("not support \"redirect\" mode of transparent proxy: %w", err)
		}
		iptables.SetWatcher(&iptables.Redirect)
	case configure.TransparentSystemProxy:
		if err = iptables.SystemProxy.GetSetupCommands().Run(true); err != nil {
			return fmt.Errorf("not support \"system proxy\" mode of transparent proxy: %w", err)
		}
	default:
		return fmt.Errorf("undefined \"%v\" mode of transparent proxy", setting.TransparentType)
	}

	if specialMode.ShouldLocalDnsListen() {
		if couldListenLocalhost, e := specialMode.CouldLocalDnsListen(); couldListenLocalhost {
			if e != nil {
				log.Warn("only listen at 127.2.0.17: %v", e)
			}
			resetResolvHijacker()
		} else if specialMode.ShouldUseFakeDns() {
			return fmt.Errorf("fakedns cannot be enabled: %w", e)
		} else {
			log.Warn("WriteTransparentProxyRules: %v", e)
		}
	}
	return nil
}

func CheckAndSetupTransparentProxy(checkRunning bool, setting *configure.Setting) (err error) {
	if setting != nil {
		setting.FillEmpty()
	} else {
		setting = configure.GetSettingNotNil()
	}
	if (!checkRunning || ProcessManager.Running()) && setting.Transparent != configure.TransparentClose {
		DeleteTransparentProxyRules()
		if !IsTransparentOn() {
			return nil
		}
		err = WriteTransparentProxyRules()
	}
	return
}

func CheckAndStopTransparentProxy() {
	DeleteTransparentProxyRules()
}

func IsTransparentOn() bool {
	setting := configure.GetSettingNotNil()
	if setting.Transparent == configure.TransparentClose {
		return false
	}
	if conf.GetEnvironmentConfig().Lite &&
		(setting.TransparentType == configure.TransparentTproxy ||
			setting.TransparentType == configure.TransparentRedirect) {
		return false
	}
	return true
}
