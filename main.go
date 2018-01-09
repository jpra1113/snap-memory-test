package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/jpra1113/snap-plugin-collector-docker/collector"
	"github.com/jpra1113/snap-plugin-lib-go/v1/plugin"
	"github.com/jpra1113/snap-plugin-publisher-influxdb/influxdb"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getPodIp(labelSelector string, namespace string) (string, error) {
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		return "", errors.New("Unable to get in cluster kubeconfig: " + err.Error())
	}

	k8sClient, err := k8s.NewForConfig(kubeConfig)
	if err != nil {
		return "", errors.New("Unable to create in cluster k8s client: " + err.Error())
	}

	pods, err := k8sClient.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return "", fmt.Errorf("Unable to list %s pod: %s", labelSelector, err.Error())
	}

	return pods.Items[0].Status.PodIP, nil
}

func main() {
	influxdb := influxdb.NewInfluxPublisher()
	influxsrvHost, err := getPodIp("app=influxsrv", "hyperpilot")
	if err != nil {
		fmt.Println("Unable to get influxsrv pod ip: " + err.Error())
	}

	influxsrvCfg := plugin.Config{
		"host":          influxsrvHost,
		"scheme":        "http",
		"port":          int64(8086),
		"user":          "root",
		"password":      "default",
		"database":      "snapaverage",
		"retention":     "autogen",
		"skip-verify":   false,
		"isMultiFields": false,
		"debug":         false,
		"log-level":     "debug",
		"precision":     "s",
	}

	docker := collector.New()
	dockerCfg := plugin.Config{
		"endpoint": "unix:///var/run/docker.sock",
		"procfs":   "/proc",
	}

	for {
		metricTypes, err := docker.GetMetricTypes(dockerCfg)
		if err != nil {
			fmt.Println("Unable to get metric types: " + err.Error())
		}

		newMetricTypes := []plugin.Metric{}
		for _, mts := range metricTypes {
			mts.Config = dockerCfg
			newMetricTypes = append(newMetricTypes, mts)
		}

		collectMetrics, err := docker.CollectMetrics(newMetricTypes)
		if err != nil {
			fmt.Println("Unable to collect metric types: " + err.Error())
		}

		// Publish
		if err := influxdb.Publish(collectMetrics, influxsrvCfg); err != nil {
			fmt.Println("Unable to publish to influxdb: " + err.Error())
		}
		fmt.Printf("Publish %d collect metrics to influxdb", len(collectMetrics))

		metricTypes = nil
		collectMetrics = nil
		newMetricTypes = nil
		err = nil

		time.Sleep(time.Second * 5)
	}
}
