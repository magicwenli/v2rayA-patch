package service

import (
	"sort"
	"time"

	"github.com/v2rayA/v2rayA/core/touch"
	"github.com/v2rayA/v2rayA/db/configure"
	"github.com/v2rayA/v2rayA/pkg/util/log"
)

func SetServer() ([]*configure.Which, error) {
	log.Info("starting update server...")
	t := touch.GenerateTouch()
	var wt []*configure.Which
	for i, v := range t.Subscriptions {
		for _, v2 := range v.Servers {
			newWt := configure.Which{TYPE: v2.TYPE, ID: v2.ID, Sub: i}
			wt = append(wt, &newWt)
		}
	}
	log.Info("stop v2ray core")
	err := StopV2ray()
	if err != nil {
		log.Fatal("%v", err)
		return nil, err
	}
	log.Info("testing http latency...")
	wt, err = TestHttpLatency(wt, 8*time.Second, 4, false)
	if err != nil {
		log.Fatal("%v", err)
		return nil, err
	}
	sort.SliceStable(wt, func(i, j int) bool {
		if len(wt[i].Latency) == len(wt[j].Latency) {
			return wt[i].Latency < wt[j].Latency
		} else {
			return len(wt[i].Latency) < len(wt[j].Latency)
		}
	})
	log.Info("setting fastest server...")
	for _, v := range t.ConnectedServers {
		err := Disconnect(configure.Which{TYPE: v.TYPE, ID: v.ID, Sub: v.Sub}, false)
		if err != nil {
			log.Fatal("%v", err)
			return nil, err
		}
	}
	for i := 0; i < 3; i++ {
		err = Connect(wt[i])
		if err != nil {
			log.Fatal("%v", err)
			return nil, err
		}
	}
	log.Info("restart v2ray core")
	err = StartV2ray()
	if err != nil {
		log.Fatal("%v", err)
		return nil, err
	}
	return wt, nil
}
