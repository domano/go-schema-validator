package main

import (
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

type simpleSchemaStore struct {
	sync.RWMutex
	store map[string]*gojsonschema.Schema
}

func NewSimpleSchemaStore() *simpleSchemaStore {
	return &simpleSchemaStore{store: make(map[string]*gojsonschema.Schema)}
}

func (sS *simpleSchemaStore) Insert(resourceName string, schema *gojsonschema.Schema) {
	sS.Lock()
	sS.store[resourceName] = schema
	defer sS.Unlock()
}

func (sS *simpleSchemaStore) Remove(resourceName string) {
	sS.Lock()
	delete(sS.store, resourceName)
	defer sS.Unlock()
}

func (sS *simpleSchemaStore) Get(resourceName string) (schema *gojsonschema.Schema, exists bool) {
	sS.RLock()
	defer sS.RUnlock()
	schema, exists = sS.store[resourceName]
	return
}
