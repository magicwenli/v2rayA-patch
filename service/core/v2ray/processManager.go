package v2ray

import (
	"fmt"
	"time"

	"github.com/magicwenli/v2rayA-patch/common/ntp"
	"github.com/magicwenli/v2rayA-patch/common/resolv"
	"github.com/magicwenli/v2rayA-patch/core/specialMode"
	"github.com/magicwenli/v2rayA-patch/core/v2ray/asset"
	"github.com/magicwenli/v2rayA-patch/db/configure"
	"github.com/magicwenli/v2rayA-patch/pkg/util/log"
)

type CoreProcessManager struct {
	p *Process
}

var ProcessManager CoreProcessManager

func (m *CoreProcessManager) beforeStop() {
	CheckAndStopTransparentProxy()
	specialMode.StopDNSSupervisor()
}

func (m *CoreProcessManager) afterStop() {
}

func (m *CoreProcessManager) Stop(saveRunning bool) {
	m.beforeStop()

	var err error
	if m.p != nil {
		err = m.p.Close()
		if err != nil {
			log.Warn("CoreProcessManager.Stop: %v", err)
		}
	}
	m.p = nil
	if saveRunning {
		configure.SetRunning(false)
	}

	m.afterStop()
}

func (m *CoreProcessManager) beforeStart(t *Template) (err error) {
	resolv.CheckResolvConf()
	if ok, t, err := ntp.IsDatetimeSynced(); err == nil && !ok {
		return fmt.Errorf("Please sync datetime first. Your datetime is %v, and the correct datetime is %v", time.Now().Local().Format(ntp.DisplayFormat), t.Local().Format(ntp.DisplayFormat))
	}

	if (t.Setting.Transparent == configure.TransparentGfwlist || t.Setting.RulePortMode == configure.GfwlistMode) && !asset.DoesV2rayAssetExist("LoyalsoldierSite.dat") {
		return fmt.Errorf("cannot find GFWList files. update GFWList and try again")
	}

	if err = t.CheckInboundPortsOccupied(); err != nil {
		return err
	}
	return nil
}

func (m *CoreProcessManager) afterStart(t *Template) (err error) {
	if err = CheckAndSetupTransparentProxy(true, t.Setting); err != nil {
		m.Stop(true)
		return
	}
	specialMode.CheckAndSetupDNSSupervisor()
	return nil
}

func (m *CoreProcessManager) Start(t *Template) (err error) {
	m.Stop(true)

	if err := m.beforeStart(t); err != nil {
		return err
	}

	process, err := NewProcess(t)
	if err != nil {
		return err
	}
	m.p = process

	if err := m.afterStart(t); err != nil {
		return err
	}

	configure.SetRunning(true)
	return nil
}

// Running reports if v2ray-core is running.
func (m *CoreProcessManager) Running() bool {
	return m.p != nil
}

func (m *CoreProcessManager) Process() *Process {
	return m.p
}
