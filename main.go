package main

import (
	"fmt"
	"net/http"

	"github.com/domano/go-schema-validator/store"

	"github.com/domano/go-schema-validator/config"
	"github.com/domano/go-schema-validator/rest"

	"github.com/gorilla/mux"

	"github.com/tarent/go-log-middleware/logging"
)

func main() {

	conf := config.NewConfig()
	//TODO: add timeoutstuff here

	server := &http.Server{Addr: fmt.Sprintf("%s:%d", "", conf.ValidationPort), ReadTimeout: conf.ValidationTimeout}

	// Init schemastore
	store := store.NewSimpleSchemaStore()

	// Init handlers and mux
	router := mux.NewRouter()
	validationHandler := rest.NewValidationHandler(store, nil)
	schemaHandler := rest.NewSchemaHandler(store, nil)

	router.Handle("/schema/{schemaName}", schemaHandler)
	router.Handle("/validation/{schemaName}", validationHandler)

	server.Handler = router
	err := server.ListenAndServe()
	if err != nil {
		logging.Logger.WithError(err).Fatal("Could not spin up http server.")
	}
}
