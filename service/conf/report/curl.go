package report

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/magicwenli/v2rayA-patch/common/httpClient"
	"github.com/magicwenli/v2rayA-patch/core/v2ray"
	"github.com/magicwenli/v2rayA-patch/db/configure"
)

type CurlReporter struct {
}

var DefaultCurlReporter CurlReporter

func (r *CurlReporter) PresetPortReport() (ok bool, report string) {
	defer func() {
		report = "Preset Port(socks5): " + report
	}()
	ports := configure.GetPortsNotNil()
	if ports.Socks5 == 0 {
		return false, "Preset HTTP Port is closed"
	}
	cli, err := httpClient.GetHttpClientWithProxy("socks5://" + net.JoinHostPort("127.0.0.1", strconv.Itoa(ports.Socks5)))
	if err != nil {
		return false, err.Error()
	}
	resp, err := cli.Get("https://www.apple.com")
	if err != nil || resp.StatusCode != 200 {
		if err == nil {
			resp.Body.Close()
			return false, resp.Status
		}
		return false, err.Error()
	}
	resp.Body.Close()
	return true, resp.Status
}

func (r *CurlReporter) TransparentReport() (ok bool, report string) {
	setting := configure.GetSettingNotNil()
	defer func() {
		report = fmt.Sprintf("Transparent Proxy(%v): %v", setting.TransparentType, report)
	}()
	if !v2ray.IsTransparentOn() {
		return true, "Transparent Proxy is not enabled"
	}
	cli := http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := cli.Get("https://ipv4.appspot.com/")
	if err != nil || resp.StatusCode != 200 {
		if err == nil {
			resp.Body.Close()
			return false, resp.Status
		}
		return false, err.Error()
	}
	b, _ := io.ReadAll(resp.Body)
	ip := strings.TrimSpace(string(b))
	if net.ParseIP(ip) == nil {
		return false, "UNKNOWN PROBLEM"
	}
	return true, "Your remote IP is: " + ip
}
