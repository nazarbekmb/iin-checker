package main

import (
	"iin_check/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.MainPageHandler).Methods("GET", "POST")
	r.HandleFunc("/iin_check/{iin}", handlers.IinCheckHandler).Methods("GET")
	r.HandleFunc("/people/info", handlers.AddPersonHandler).Methods("GET", "POST")
	r.HandleFunc("/people/info/iin/{iin}", handlers.GetPersonByIINHandler).Methods("GET")
	r.HandleFunc("/people/info/phone/{name}", handlers.GetPersonsByNameHandler).Methods("GET")

	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	log.Println("Starting server on port 8081...http://localhost:8081/")
	log.Fatal(http.ListenAndServe(":8081", r))
}
