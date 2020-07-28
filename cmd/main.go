package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

type LogMsg struct {
	IP string `json:"ip"`
}

var (
	uniqueIPAddresses = make(map[string]struct{})
	rwMutex = sync.RWMutex{}
)

func main() {

	log := logrus.New()

	cfg, err := NewConfig()
	if err != nil {
		log.WithError(err).Fatalln("could not parse config")
	}

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.WithField("logLevel", cfg.LogLevel).WithError(err).Errorln("could not parse log level")
	}

	log.SetLevel(logLevel)

	registry := prometheus.NewRegistry()

	ipCollector := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "unique_ip_addresses",
		Help:        "Count of unique ip address.",
	}, func() float64 {

		defer rwMutex.RUnlock()
		rwMutex.RLock()

		return float64(len(uniqueIPAddresses))
	})

	registry.MustRegister(ipCollector)

	metricsHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	metricsRouter := mux.NewRouter()
	metricsRouter.Handle("/metrics", metricsHandler)

	metricsSrv := http.Server{
		Handler:      metricsRouter,
		Addr:         fmt.Sprintf(":%s", cfg.MetricsPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		err := metricsSrv.ListenAndServe()
		if err != nil {
			log.WithError(err).Fatalln("could not start metrics server")
		}
	}()

	apiRouter := mux.NewRouter()
	apiRouter.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(200)

		logMsg := &LogMsg{}
		err := json.NewDecoder(r.Body).Decode(logMsg)
		if err != nil {
			log.WithError(err).Errorln("could not parse body")
			return
		}

		defer rwMutex.Unlock()
		rwMutex.Lock()

		uniqueIPAddresses[logMsg.IP] = struct{}{}
	}).Methods("POST").Headers("Content-Type", "application/json")

	apiSrv := http.Server{
		Handler:      apiRouter,
		Addr:         fmt.Sprintf(":%s", cfg.ApiPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err = apiSrv.ListenAndServe()
	if err != nil {
		log.WithError(err).Fatalln("could not start api server")
	}

}
