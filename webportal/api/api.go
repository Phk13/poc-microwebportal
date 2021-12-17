package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/phk13/poc-micro/databaselayer"
)

func RunApi(endpoint string, db databaselayer.DinoDBHandler) error {
	r := mux.NewRouter()
	RunApiOnRouter(r, db)
	return http.ListenAndServe(endpoint, r)
}

func RunApiOnRouter(r *mux.Router, db databaselayer.DinoDBHandler) {
	handler := newDinoRESTAPIHandler(db)

	apirouter := r.PathPrefix("/api/dinos").Subrouter()
	apirouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.searchHandler)
	apirouter.Methods("POST").PathPrefix("/{Operation}").HandlerFunc(handler.editsHandler)
}
