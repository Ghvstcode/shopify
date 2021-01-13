package main

import (
	"context"
	l "github.com/GhvstCode/shopify-challenge/utils/logger"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	//utils.DbCron()
}


func handleRequest() {
	r := mux.NewRouter().StrictSlash(true)
	u := r.PathPrefix("/api/v1").Subrouter()




	s := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      r,
		IdleTimeout:  5 * time.Minute, //120
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

	go func() {
		l.InfoLogger.Println("Server is up on port", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			l.ErrorLogger.Println(err)
			l.ErrorLogger.Fatal("Error starting server on port", s.Addr)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.InfoLogger.Println("Received terminate, graceful shutdown! Signal: ", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_ = s.Shutdown(tc)

}