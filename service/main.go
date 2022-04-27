package main

import (
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/magicwenli/v2rayA-patch/conf/report"
	_ "github.com/magicwenli/v2rayA-patch/pkg/plugin/pingtunnel"
	_ "github.com/magicwenli/v2rayA-patch/pkg/plugin/simpleobfs"
	_ "github.com/magicwenli/v2rayA-patch/pkg/plugin/socks5"
	_ "github.com/magicwenli/v2rayA-patch/pkg/plugin/ss"
	_ "github.com/magicwenli/v2rayA-patch/pkg/plugin/ssr"
	_ "github.com/magicwenli/v2rayA-patch/pkg/plugin/tcp"
	_ "github.com/magicwenli/v2rayA-patch/pkg/plugin/tls"
	_ "github.com/magicwenli/v2rayA-patch/pkg/plugin/trojanc"
	_ "github.com/magicwenli/v2rayA-patch/pkg/plugin/ws"
	"github.com/magicwenli/v2rayA-patch/pkg/util/log"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	checkEnvironment()
	if runtime.GOOS == "linux" {
		checkTProxySupportability()
	}
	initConfigure()
	checkUpdate()
	hello()
	if err := run(); err != nil {
		log.Fatal("main: %v", err)
	}
}
