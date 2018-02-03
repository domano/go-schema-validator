package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/tarent/go-log-middleware/logging"
	"github.com/xeipuuv/gojsonschema"
)

type validationHandler struct {
	*simpleSchemaStore
	next http.Handler
}

func NewValidationHandler(schemaStr string, delegate http.Handler) (*validationHandler, error) {
	vH := &validationHandler{simpleSchemaStore: NewSimpleSchemaStore(), next: delegate}
	err := vH.SetSchema("product", schemaStr)
	if err != nil {
		return nil, errors.Wrap(err, "Could not set schema for new validation handler: ")
	}
	return vH, err
}

func (vH *validationHandler) SetSchema(name, schemaStr string) error {
	loader := gojsonschema.NewStringLoader(schemaStr)
	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return errors.Wrap(err, "Could not set json schema for validationhandler ")
	}
	vH.Insert(name, schema)
	return nil
}

func (vH *validationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	jsonLoader, newBody := gojsonschema.NewReaderLoader(r.Body)

	//TODO use gorilla and parse path params here
	schemaName := strings.ToLower(strings.Replace(r.URL.EscapedPath(), "/", "", -1))
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
			stringBuf.WriteString(validationErrors[i].String())
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

}
