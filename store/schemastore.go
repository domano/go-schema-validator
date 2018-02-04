package store

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

type SimpleSchemaStore struct {
	sync.RWMutex
	store     map[string]*gojsonschema.Schema
	byteStore map[string][]byte
}

func NewSimpleSchemaStore() *SimpleSchemaStore {
	return &SimpleSchemaStore{store: make(map[string]*gojsonschema.Schema), byteStore: make(map[string][]byte)}
}

func (sS *SimpleSchemaStore) Insert(resourceName string, schemaStr []byte) error {
	sS.Lock()
	defer sS.Unlock()
	loader := gojsonschema.NewBytesLoader(schemaStr)
	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return errors.Wrap(err, "Could not set json schema for schemaHandler ")
	}
	sS.store[resourceName] = schema
	sS.byteStore[resourceName] = schemaStr
	return nil
}

func (sS *SimpleSchemaStore) Remove(resourceName string) {
	sS.Lock()
	defer sS.Unlock()
	delete(sS.store, resourceName)
}

func (sS *SimpleSchemaStore) Get(resourceName string) (schema *gojsonschema.Schema, exists bool) {
	sS.RLock()
	defer sS.RUnlock()
	schema, exists = sS.store[resourceName]
	return
}

func (sS *SimpleSchemaStore) GetBytes(schemaName string) (schemaBytes []byte, exists bool) {
	sS.RLock()
	defer sS.RUnlock()
	schemaBytes, exists = sS.byteStore[schemaName]
	return
}
