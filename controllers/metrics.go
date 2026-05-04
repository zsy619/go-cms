package controllers

import (
	"net/http/pprof"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ProfController struct {
	web.Controller
}

func (this *ProfController) Get() {
	switch this.Ctx.Input.Param(":pp") {
	default:
		pprof.Index(this.Ctx.ResponseWriter, this.Ctx.Request)
	case "":
		pprof.Index(this.Ctx.ResponseWriter, this.Ctx.Request)
	case "cmdline":
		pprof.Cmdline(this.Ctx.ResponseWriter, this.Ctx.Request)
	case "profile":
		pprof.Profile(this.Ctx.ResponseWriter, this.Ctx.Request)
	case "symbol":
		pprof.Symbol(this.Ctx.ResponseWriter, this.Ctx.Request)
	}
	this.Ctx.ResponseWriter.WriteHeader(200)
}

func recordMetrics() {
	// 每2秒，计数器增加1。
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
	Name: "myapp_processed_ops_total",
	Help: "The total number of processed events",
})

func init() {
	recordMetrics()
	web.Handler("/debug/metrics", promhttp.Handler())

	web.Router(`/debug/:pp([\w]+)`, &ProfController{})
}
