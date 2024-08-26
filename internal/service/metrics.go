package service

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	echo "github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/slyngshot-al/packages/log"
)

var initialized bool
var actionCount *prometheus.CounterVec
var actionDuration *prometheus.SummaryVec

var regOnce = sync.Once{}

func InitMetrics(namespace string, addr string) (err error) {
	regOnce.Do(func() {
		actionCount = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Name: "action_count"}, []string{"action", "state"})
		if err = prometheus.DefaultRegisterer.Register(actionCount); err != nil {
			return
		}

		actionDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{Namespace: namespace, Name: "action_duration_ms"}, []string{"action", "state"})
		if err = prometheus.DefaultRegisterer.Register(actionDuration); err != nil {
			return
		}

		initialized = true
		go func() {
			metrics := echo.New()                                // this Echo will run on separate port 9090
			metrics.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
			if srvErr := metrics.Start(addr); srvErr != nil && !errors.Is(srvErr, http.ErrServerClosed) {
				log.Fatal(context.Background(), srvErr, "metrics server failed")
			}
		}()
	})

	return err
}

func CollectMetricFn(action string) func(ctx context.Context, err error) {
	startTime := time.Now()
	return func(ctx context.Context, err error) {
		state := "ok"
		if err != nil {
			log.Error(ctx, err, "action failed")
			state = "failed"
		}
		if !initialized {
			return
		}
		actionCount.WithLabelValues(action, state).Inc()
		actionDuration.WithLabelValues(action, state).Observe(float64(time.Since(startTime).Milliseconds()))
	}
}
