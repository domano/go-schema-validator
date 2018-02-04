package rest

import (
	"bytes"
	"fmt"
	"go-schema-validator/store"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"

	"github.com/tarent/go-log-middleware/logging"
)

type validationHandler struct {
	*store.SimpleSchemaStore
	next http.Handler
}

func NewValidationHandler(store *store.SimpleSchemaStore, delegate http.Handler) *validationHandler {
	return &validationHandler{SimpleSchemaStore: store, next: delegate}
}

func (vH *validationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	jsonLoader, tReader := gojsonschema.NewReaderLoader(r.Body)
	newBody := &bytes.Buffer{}
	_, err := io.Copy(newBody, tReader)
	if err != nil && err != io.EOF {
		logging.Logger.WithError(err).Error("Could not read body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	pathParams := mux.Vars(r)
	schemaName := pathParams["schemaName"]
	schema, exists := vH.Get(schemaName)
	if !exists {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	res, err := schema.Validate(jsonLoader)
	if err != nil && err != io.EOF {
		logging.Logger.WithError(err).Error("Could not parse body as json")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if !res.Valid() {
		validationErrors := res.Errors()
		stringBuf := bytes.NewBuffer(make([]byte, 0))
		for i := range validationErrors {
			stringBuf.WriteString(fmt.Sprintln(validationErrors[i].String()))
		}
		_, err := rw.Write(stringBuf.Bytes())
		if err != nil {
			// Could not write to respWriter
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	if vH.next != nil {
		r.Body = ioutil.NopCloser(newBody)
		vH.next.ServeHTTP(rw, r)
	}
	return

}
