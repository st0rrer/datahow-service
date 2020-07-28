package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Service interface {
	CountUniqueIP() (int, error)
}

func NewMetricHandler(service Service) http.Handler {

	registry := prometheus.NewRegistry()

	ipCollector := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "unique_ip_addresses",
		Help:        "Count of unique ip address.",
	}, func() float64 {

		count, err := service.CountUniqueIP()
		if err != nil {
			return 0
		}

		return float64(count)
	})

	registry.MustRegister(ipCollector)

	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}
