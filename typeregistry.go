// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package typeregistry implements a simple type registry that maps go types to custom names.
package typeregistry

import (
	"reflect"
	"sync"

	"github.com/vedranvuk/errorex"
)

var (
	// ErrTypeRegistry is the base typeregistry error.
	ErrTypeRegistry = errorex.New("typeregistry")
	// ErrNotFound is returned when a type was not found in the registry.
	ErrNotFound = ErrTypeRegistry.WrapFormat("entry '%s' not found")
	// ErrDuplicateEntry is returned when registering a type that is alreday registered.
	ErrDuplicateEntry = ErrTypeRegistry.WrapFormat("entry '%s' already exists")
)

// Registry is a Go type registry.
type Registry struct {
	mu      sync.Mutex
	entries map[string]reflect.Type
}

// New creates a new Registry instance.
func New() *Registry {
	p := &Registry{
		mu:      sync.Mutex{},
		entries: make(map[string]reflect.Type),
	}
	return p
}

// Register registers a type of value v under specified name.
func (r *Registry) Register(name string, v interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.entries[name]; exists {
		return ErrDuplicateEntry.WrapArgs(name)
	}

	r.entries[name] = reflect.TypeOf(v)

	return nil
}

// Unregister unregisters a type by name/type name.
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.entries[name]; !exists {
		return ErrNotFound.WrapArgs(name)
	}
	delete(r.entries, name)

	return nil
}

// Get creates a new instance of a registered type by specified name.
func (r *Registry) Get(name string) (interface{}, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.entries[name]
	if !ok {
		return nil, ErrNotFound.WrapArgs(name)
	}
	return reflect.New(t).Elem().Interface(), nil
}

// RegisteredNames returns a slice of registered names.
func (r *Registry) RegisteredNames() []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	names := make([]string, 0, len(r.entries))
	for key := range r.entries {
		names = append(names, key)
	}
	return names
}
