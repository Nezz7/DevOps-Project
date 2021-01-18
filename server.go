package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
)

func main() {

	fmt.Printf("Starting server at port 8080\n")
	route := mux.NewRouter()
	
	s := route.PathPrefix("/api").Subrouter()
	s.HandleFunc("/user", createProfile).Methods("POST")
	s.HandleFunc("/users", getAllUsers).Methods("GET")
	s.Path("/metrics").Handler(promhttp.Handler())
	log.Println(http.ListenAndServe(":8080", s))

}
