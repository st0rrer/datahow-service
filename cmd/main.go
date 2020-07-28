package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/st0rrer/datahow-service/src/log"
	"github.com/st0rrer/datahow-service/src/metrics"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {

	logger := logrus.New()

	cfg, err := NewConfig()
	if err != nil {
		logger.WithError(err).Fatalln("could not parse config")
	}

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.WithField("logLevel", cfg.LogLevel).WithError(err).Errorln("could not parse logger level")
	}

	logger.SetLevel(logLevel)

	logService := log.NewService()
	apiHandler := log.Handler{Service: logService}

	apiRouter := mux.NewRouter()
	apiRouter.HandleFunc("/logs", apiHandler.ProcessMessage).Methods("POST").Headers("Content-Type", "application/json")

	apiSrv := http.Server{
		Handler:      apiRouter,
		Addr:         fmt.Sprintf(":%s", cfg.ApiPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		err := apiSrv.ListenAndServe()
		if err != nil {
			logger.WithError(err).Fatalln("could not start api server")
		}
	}()

	metricsRouter := mux.NewRouter()
	metricsRouter.Handle("/metrics", metrics.NewMetricHandler(logService))

	metricsSrv := http.Server{
		Handler:      metricsRouter,
		Addr:         fmt.Sprintf(":%s", cfg.MetricsPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		err := metricsSrv.ListenAndServe()
		if err != nil {
			logger.WithError(err).Fatalln("could not start metrics server")
		}
	}()


	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt)

	<-signalCh

	wg := sync.WaitGroup{}

	wg.Add(2)

	go shutDownSrv(&wg, &metricsSrv)
	go shutDownSrv(&wg, &apiSrv)

	wg.Wait()
	logger.Infoln("server is successfully shutdown")
}

func shutDownSrv(wg *sync.WaitGroup, srv *http.Server) {

	defer wg.Done()

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	srv.Shutdown(ctx)
}
