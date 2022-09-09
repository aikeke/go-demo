package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
	"zabbix-alert-tool/monitor"
)

func main() {

	cnHttp := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cn_http_monitor",
			Help: "cn http interface status",
		},
		// 指定标签名称
		[]string{"area", "metrics", "msg"},
	)
	cnPrivate := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cn_private_monitor",
			Help: "cn private interface status",
		},
		// 指定标签名称
		[]string{"area", "metrics", "msg"},
	)
	cnWebsocket := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cn_websocket_monitor",
			Help: "cn websocket interface status",
		},
		// 指定标签名称
		[]string{"area", "metrics", "msg"},
	)
	enHttp := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "en_http_monitor",
			Help: "en http interface status",
		},
		// 指定标签名称
		[]string{"area", "metrics", "msg"},
	)
	enPrivate := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "en_private_monitor",
			Help: "en private interface status",
		},
		// 指定标签名称
		[]string{"area", "metrics", "msg"},
	)
	enWebsocket := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "en_websocket_monitor",
			Help: "en websocket interface status",
		},
		// 指定标签名称
		[]string{"area", "metrics", "msg"},
	)

	// 注册到全局默认注册表中
	prometheus.MustRegister(cnHttp, cnPrivate, cnWebsocket)
	prometheus.MustRegister(enHttp, enPrivate, enWebsocket)
	go func() {
		for {
			code, msg := monitor.CnHttpMonitor("广州", "http://cn-edugz.hivoice.cn/eval/mp3")
			cnHttp.WithLabelValues("gz", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.CnHttpMonitor("上海", "http://cn-edush.hivoice.cn/eval/mp3")
			cnHttp.WithLabelValues("sh", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.CnHttpMonitor("北京", "http://cn-edubj.hivoice.cn/eval/mp3")
			cnHttp.WithLabelValues("bj", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.EnHttpMonitor("广州", "http://edugz.hivoice.cn/eval/mp3")
			cnHttp.WithLabelValues("gz", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.EnHttpMonitor("上海", "http://edush.hivoice.cn/eval/mp3")
			cnHttp.WithLabelValues("sh", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.EnHttpMonitor("北京", "http://edubj.hivoice.cn/eval/mp3")
			cnHttp.WithLabelValues("bj", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			//private
			code, msg = monitor.CnPrivMonitor("广州", "cn-evalgz.hivoice.cn:18085")
			cnPrivate.WithLabelValues("gz", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.CnPrivMonitor("上海", "cn-evalsh.hivoice.cn:18085")
			cnPrivate.WithLabelValues("sh", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.CnPrivMonitor("北京", "cn-evalbj.hivoice.cn:18085")
			cnPrivate.WithLabelValues("bj", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.EnPrivMonitor("广州", "esg.hivoice.cn:8085")
			enPrivate.WithLabelValues("gz", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.EnPrivMonitor("上海", "est.hivoice.cn:8085")
			enPrivate.WithLabelValues("sh", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.EnPrivMonitor("北京", "esb.hivoice.cn:8085")
			enPrivate.WithLabelValues("bj", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			//websocket
			code, msg = monitor.CnWsMonitor("广州", "172.18.8.220:18081")
			cnWebsocket.WithLabelValues("gz", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.CnWsMonitor("上海", "172.17.8.220:18081")
			cnWebsocket.WithLabelValues("sh", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.CnWsMonitor("北京", "172.16.8.220:18081")
			cnWebsocket.WithLabelValues("bj", "code", msg).Set(float64(code))
			code, msg = monitor.EnWsMonitor("广州", "172.18.8.220:8081")
			enWebsocket.WithLabelValues("gz", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.EnWsMonitor("上海", "172.17.8.220:8081")
			enWebsocket.WithLabelValues("sh", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 1)
			code, msg = monitor.EnWsMonitor("北京", "172.16.8.220:8081")
			enWebsocket.WithLabelValues("bj", "code", msg).Set(float64(code))
			time.Sleep(time.Second * 180)
		}
	}()

	// 暴露自定义的指标
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
