package rest

import (
	"net/http"

	"EventPlatform/myevent/src/lib/persistence"
	"github.com/gorilla/mux"
)

func ServeAPI(endpoint string, databasehandler persistence.DatabaseHandler) error {
	handler := NewEventHandler(databasehandler)
	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()
	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.FindEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.AllEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.NewEventHandler)
	return http.ListenAndServe(endpoint, r)
}
