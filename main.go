package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	rg := run.Group{}

	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "demo_count",
		Help: "Demo counter",
	})

	prometheus.MustRegister(counter)

	// Every second add a random value [0..10) to the counter.
	go func() {
		for {
			val := rand.Int63n(10)
			counter.Add(float64(val))

			time.Sleep(time.Second)
		}
	}()

	// Handle interrupts from the OS for graceful shutdown
	{
		// Catch SIGINT (ctrl-c) and SIGTERM to gracefully shutdown. SIGTERM is
		// sent by Heroku when restarting dynos.
		interrupts := make(chan os.Signal, 1)
		signal.Notify(interrupts, syscall.SIGINT, syscall.SIGTERM)

		rg.Add(
			func() error { <-interrupts; return nil },
			func(err error) {},
		)
	}

	{
		prom := prometheus.Handler()

		http.HandleFunc("/metrics", func(resp http.ResponseWriter, req *http.Request) {
			log.Printf("%s %s\n", req.Method, req.URL.Path)
			prom.ServeHTTP(resp, req)
		})

		server := http.Server{
			Addr:    ":1845",
			Handler: http.DefaultServeMux,
		}

		rg.Add(
			func() error {
				err := server.ListenAndServe()

				// Ignore the error returned when the server is closed manually.
				if err == http.ErrServerClosed {
					err = nil
				}

				return err
			},
			func(err error) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				// Gracefully shutdown the HTTP listener.
				err = server.Shutdown(ctx)
				if err != nil {
					log.Println(err)
				}
			},
		)
	}

	err := rg.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
