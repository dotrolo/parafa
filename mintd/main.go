package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dotrolo/parafa/mintd/internal/admin"
	"github.com/dotrolo/parafa/mintd/internal/api"
)

func main() {
	api := &http.Server{
		Addr:         "127.0.0.1:8080",
		Handler:      api.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	admin := &http.Server{
		Addr:         "127.0.0.1:8081",
		Handler:      admin.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errs := make(chan error, 2)

	go func() {
		log.Println("api server: " + api.Addr)
		errs <- api.ListenAndServe()
	}()

	go func() {
		log.Println("admin server: " + admin.Addr)
		errs <- admin.ListenAndServe()
	}()

	log.Fatal(<-errs)
}
