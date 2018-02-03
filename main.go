package main

import (
	"net/http"

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
	validationServer := &http.Server{Addr: ":9999"}

	validationHandler, err := NewValidationHandler(schema, nil)
	if err != nil {
		logging.Logger.WithError(err).Fatal("Could not read intial schema for validationhandler.")
	}
	validationServer.Handler = validationHandler
	err = validationServer.ListenAndServe()
	if err != nil {
		logging.Logger.WithError(err).Fatal("Could not spin up http server.")
	}
}
