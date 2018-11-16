package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	//
	"negorni_test/middleware"
)

var (
	osSignal         chan os.Signal
	globalMiddleware middleware.GlobalMiddleware
	groupMiddleware  middleware.GroupMiddleware
)

func init() {
	osSignal = make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	// init middleware
	globalMiddleware = middleware.GlobalMiddle
	groupMiddleware = middleware.GroupMiddle
}

func main() {
	negro := negroni.New()
	routes := mux.NewRouter()

	// ex. for route group middleware
	subRouter := mux.NewRouter().PathPrefix("/group").Subrouter().StrictSlash(true)
	subRouter.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to group!")
	})

	routes.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to home!")
	})
	routes.PathPrefix("/group").Handler(negroni.New(groupMiddleware, negroni.Wrap(subRouter)))

	// ex. for global middleware
	negro.Use(globalMiddleware.Validate())
	negro.UseHandler(routes)

	srv := &http.Server{
		Addr: "0.0.0.0:9090",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      negro, // Pass our instance of gorilla/mux in.
	}

	log.Printf("Server is started\n")

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Start server error : %s\n", err)
			os.Exit(0)
		}

	}()
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	select {
	case <-osSignal:
		log.Printf("Server Shutting down\n")
		time.Sleep(time.Second * 2)
		srv.Shutdown(ctx)
		os.Exit(0)
	}
}
