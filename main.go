package main

import (
	"go-schema-validator/rest"
	"go-schema-validator/store"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tarent/go-log-middleware/logging"
)

var schema = `{
    "$schema": "http://json-schema.org/draft-06/schema#",
    "title": "Product",
    "description": "A product from Acme's catalog",
    "type": "object",
    "properties": {
        "id": {
            "description": "The unique identifier for a product",
            "type": "integer"
        },
        "name": {
            "description": "Name of the product",
            "type": "string"
        }
    },
    "required": ["id", "name"]
}`

func main() {

	//TODO: add timeoutstuff here
	server := &http.Server{Addr: ":9999"}

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
