package rest

import (
	"io"
	"io/ioutil"
	"net/http"

	"go-schema-validator/store"

	"github.com/gorilla/mux"
	"github.com/tarent/go-log-middleware/logging"
)

type schemaHandler struct {
	*store.SimpleSchemaStore
	next http.Handler
}

func NewSchemaHandler(store *store.SimpleSchemaStore, delegate http.Handler) *schemaHandler {
	return &schemaHandler{SimpleSchemaStore: store, next: delegate}
}

func (sH *schemaHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		pathParams := mux.Vars(r)
		schemaName := pathParams["schemaName"]
		schemaBytes, err := ioutil.ReadAll(r.Body)
		if err != nil && err != io.EOF {
			logging.Logger.WithError(err).Info("Could not parse request body for schemahandler")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		err = sH.Insert(schemaName, schemaBytes)
		if err != nil{
			logging.Logger.WithError(err).Info("Could not save schema in schemahandler")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		return

	case http.MethodDelete:
		pathParams := mux.Vars(r)
		schemaName := pathParams["schemaName"]
		sH.Remove(schemaName)
		return

	case http.MethodGet:
		pathParams := mux.Vars(r)
		schemaName := pathParams["schemaName"]
		schema, exists := sH.GetBytes(schemaName)
		if !exists {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		_, err := rw.Write(schema)
		if err != nil {
			logging.Logger.WithError(err).Error("Could not respond with schema string")
			rw.WriteHeader(http.StatusInternalServerError)
		}
		return
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
