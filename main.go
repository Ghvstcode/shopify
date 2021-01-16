package main

import (
	"context"
	"github.com/GhvstCode/shopify-challenge/controllers"
	"github.com/GhvstCode/shopify-challenge/middleware"
	l "github.com/GhvstCode/shopify-challenge/utils/logger"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	//utils.DbCron()
	handleRequest()
}


func handleRequest() {
	r := mux.NewRouter().StrictSlash(true)
	u := r.PathPrefix("/api/v1").Subrouter()


	//r.Use(MiddleWare.Jwt2)
	//Register MiddleWare.
	r.Use(middleware.Jwt)

	//r.HandleFunc("/", controllers.Home).Methods(http.MethodGet)
	//r.HandleFunc("/logs", controllers.ViewLog).Methods(http.MethodGet)
	u.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)
	u.HandleFunc("/signup", controllers.SignUp).Methods(http.MethodPost)
	u.HandleFunc("/upload", controllers.UploadImage).Methods(http.MethodPost)


	s := &http.Server{
		//Addr:         ":" + os.Getenv("PORT"),
		Addr:         ":3000",
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