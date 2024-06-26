package app

import (
	"context"
	"flag"
	"log"
	"multithreaded-data-processor/internal/handler"
	"multithreaded-data-processor/internal/resultData"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
)

func Run() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: service(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	sigs := make(chan os.Signal)
	signal.Notify(sigs,
		syscall.SIGINT,
		os.Interrupt)

	<-sigs
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

func service() http.Handler {
	storage := resultData.NewStorage()
	router := chi.NewRouter()
	handler.Build(router, storage)

	return router
}
