package main

import (
	"fmt"
	"time"

	"github.com/jpra1113/snap-plugin-collector-docker/collector"
	"github.com/jpra1113/snap-plugin-lib-go/v1/plugin"
)

func main() {
	for {
		cfg := plugin.Config{
			"endpoint": "unix:///var/run/docker.sock",
			"procfs":   "/proc",
		}
		docker := collector.New()
		metricTypes, err := docker.GetMetricTypes(cfg)
		if err != nil {
			fmt.Println("Unable to get metric types: " + err.Error())
		}

		newMetricTypes := []plugin.Metric{}
		for _, mts := range metricTypes {
			mts.Config = cfg
			newMetricTypes = append(newMetricTypes, mts)
		}

		collectMetrics, err := docker.CollectMetrics(newMetricTypes)
		if err != nil {
			fmt.Println("Unable to collect metric types: " + err.Error())
		}
		fmt.Println(collectMetrics)

		time.Sleep(time.Second * 5)
	}
}
