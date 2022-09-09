package main

import (
	// "fmt"
	"fmt"
	"net/http"
	"sync"
	"time"
	"zabbix-alert-tool/monitor"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var cnHttpMap = make(map[string]string)
var cnWebsocketMap = make(map[string]string)
var cnPrivateMap = make(map[string]string)
var enHttpMap = make(map[string]string)
var enWebsocketMap = make(map[string]string)
var enPrivateMap = make(map[string]string)

var wg sync.WaitGroup

type MonitorManager struct {
	CnhttpDesc      *prometheus.Desc
	EnhttpDesc      *prometheus.Desc
	CnWebsocketDesc *prometheus.Desc
	EnWebsocketDesc *prometheus.Desc
	CnPrivateDesc   *prometheus.Desc
	EnPrivateDesc   *prometheus.Desc
	// ... many more fields
}

// Simulate prepare the data
func Initmap() {
	// Just example fake data.

	cnHttpMap["广州"] = "http://cn-edugz.hivoice.cn/eval/mp3"
	cnHttpMap["上海"] = "http://cn-edush.hivoice.cn/eval/mp3"
	cnHttpMap["北京"] = "http://cn-edubj.hivoice.cn/eval/mp3"

	cnWebsocketMap["广州"] = "172.18.8.220:18081"
	cnWebsocketMap["上海"] = "172.17.8.220:18081"
	cnWebsocketMap["北京"] = "172.16.8.220:18081"

	cnPrivateMap["广州"] = "cn-evalgz.hivoice.cn:18085"
	cnPrivateMap["上海"] = "cn-evalsh.hivoice.cn:18085"
	cnPrivateMap["北京"] = "cn-evalbj.hivoice.cn:18085"

	enHttpMap["广州"] = "http://edugz.hivoice.cn/eval/mp3"
	enHttpMap["上海"] = "http://edush.hivoice.cn/eval/mp3"
	enHttpMap["北京"] = "http://edubj.hivoice.cn/eval/mp3"

	enWebsocketMap["广州"] = "172.18.8.220:8081"
	enWebsocketMap["上海"] = "172.17.8.220:8081"
	enWebsocketMap["北京"] = "172.16.8.220:8081"

	enPrivateMap["广州"] = "esg.hivoice.cn:8085"
	enPrivateMap["上海"] = "est.hivoice.cn:8085"
	enPrivateMap["北京"] = "esb.hivoice.cn:8085"
	return
}

// Describe simply sends the two Descs in the struct to the channel.
func (c *MonitorManager) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.EnhttpDesc
	ch <- c.CnhttpDesc
	ch <- c.CnPrivateDesc
	ch <- c.CnWebsocketDesc
	ch <- c.EnPrivateDesc
	ch <- c.EnWebsocketDesc
}

func (c *MonitorManager) Collect(ch chan<- prometheus.Metric) {
	// oomCountByHost, ramUsageByHost := c.ReallyExpensiveAssessmentOfTheSystemState()
	wg.Add(6)
	go func() {
		defer wg.Done()
		for area, url := range cnHttpMap {
			code, msg := monitor.CnHttpMonitor(area, url)
			ch <- prometheus.MustNewConstMetric(
				c.CnhttpDesc,
				prometheus.GaugeValue,
				float64(code),
				area,
				msg,
			)
			time.Sleep(time.Second * 1)
		}
		fmt.Println("cnhttp done")
	}()
	go func() {
		defer wg.Done()
		for area, url := range cnWebsocketMap {
			code, msg := monitor.CnWsMonitor(area, url)
			ch <- prometheus.MustNewConstMetric(
				c.CnWebsocketDesc,
				prometheus.GaugeValue,
				float64(code),
				area,
				msg,
			)
			time.Sleep(time.Second * 1)
		}
		fmt.Println("cnws done")
	}()
	go func() {
		defer wg.Done()
		for area, url := range cnPrivateMap {
			code, msg := monitor.CnPrivMonitor(area, url)
			ch <- prometheus.MustNewConstMetric(
				c.CnPrivateDesc,
				prometheus.GaugeValue,
				float64(code),
				area,
				msg,
			)
			time.Sleep(time.Second * 1)
		}
		fmt.Println("cnprivate done")
	}()
	go func() {
		defer wg.Done()
		for area, url := range enHttpMap {
			code, msg := monitor.EnHttpMonitor(area, url)
			ch <- prometheus.MustNewConstMetric(
				c.EnhttpDesc,
				prometheus.GaugeValue,
				float64(code),
				area,
				msg,
			)
			time.Sleep(time.Second * 1)
		}
		fmt.Println("enhttp done")
	}()
	go func() {
		defer wg.Done()
		for area, url := range enWebsocketMap {
			code, msg := monitor.EnWsMonitor(area, url)
			ch <- prometheus.MustNewConstMetric(
				c.EnWebsocketDesc,
				prometheus.GaugeValue,
				float64(code),
				area,
				msg,
			)
			time.Sleep(time.Second * 1)
		}
		fmt.Println("enws done")
	}()
	go func() {
		defer wg.Done()
		for area, url := range enPrivateMap {
			code, msg := monitor.EnPrivMonitor(area, url)
			ch <- prometheus.MustNewConstMetric(
				c.EnPrivateDesc,
				prometheus.GaugeValue,
				float64(code),
				area,
				msg,
			)
			time.Sleep(time.Second * 1)

		}
		fmt.Println("enprivate done")
	}()
	wg.Wait()
}

// NewMonitorManager creates the two Descs OOMCountDesc and RAMUsageDesc. Note
// that the zone is set as a ConstLabel. (It's different in each instance of the
// MonitorManager, but constant over the lifetime of an instance.) Then there is
// a variable label "host", since we want to partition the collected metrics by
// host. Since all Descs created in this way are consistent across instances,
// with a guaranteed distinction by the "zone" label, we can register different
// MonitorManager instances with the same registry.
func NewMonitorManager() *MonitorManager {
	return &MonitorManager{

		CnhttpDesc: prometheus.NewDesc(
			"cn_http_monitor",
			"cn http interface status",
			[]string{"area", "msg"},
			prometheus.Labels{},
		),
		EnhttpDesc: prometheus.NewDesc(
			"en_http_monitor",
			"en http interface status",
			[]string{"area", "msg"},
			prometheus.Labels{},
		),
		CnPrivateDesc: prometheus.NewDesc(
			"cn_private_monitor",
			"cn private interface status",
			[]string{"area", "msg"},
			prometheus.Labels{},
		),
		EnPrivateDesc: prometheus.NewDesc(
			"en_private_monitor",
			"en private interface status",
			[]string{"area", "msg"},
			prometheus.Labels{},
		),
		CnWebsocketDesc: prometheus.NewDesc(
			"cn_websocket_monitor",
			"cn websocket interface status",
			[]string{"area", "msg"},
			prometheus.Labels{},
		),
		EnWebsocketDesc: prometheus.NewDesc(
			"en_websocket_monitor",
			"en websocket interface status",
			[]string{"area", "msg"},
			prometheus.Labels{},
		),
	}
}

func main() {
	workerDB := NewMonitorManager()
	Initmap()

	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(workerDB)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":8888", nil)
}
