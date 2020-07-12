// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package typeregistry implements a simple type registry.
package typeregistry

import (
	"reflect"
	"sort"
	"sync"

	"github.com/vedranvuk/errorex"
)

var (
	// ErrTypeRegistry is the base typeregistry error.
	ErrTypeRegistry = errorex.New("typeregistry")
	// ErrNotFound is returned when a type was not found in the registry.
	ErrNotFound = ErrTypeRegistry.WrapFormat("entry '%s' not found")
	// ErrDuplicateEntry is returned when registering a type that is already
	// registered.
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

// Register registers a reflect.Type of value specified by v under its'
// reflect.Type name or returns an error.
func (r *Registry) Register(v interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t := reflect.ValueOf(v).Type()
	if _, exists := r.entries[t.String()]; exists {
		return ErrDuplicateEntry.WrapArgs(t.String())
	}

	r.entries[t.String()] = t

	return nil
}

// RegisterNamed registers a reflect.Type of value specified by v under
// specified name or returns an error.
func (r *Registry) RegisterNamed(name string, v interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.entries[name]; exists {
		return ErrDuplicateEntry.WrapArgs(name)
	}

	r.entries[name] = reflect.TypeOf(v)

	return nil
}

// Unregister unregisters a reflect.Type registered under specified name or
// returns an error.
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.entries[name]; !exists {
		return ErrNotFound.WrapArgs(name)
	}
	delete(r.entries, name)

	return nil
}

// GetType returns a registered reflect.Type specified by name or an error.
func (r *Registry) GetType(name string) (reflect.Type, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.entries[name]
	if !ok {
		return nil, ErrNotFound.WrapArgs(name)
	}
	return t, nil
}

// GetValue returns a new reflect.Value of reflect.Type registered under
// specified name or an error.
func (r *Registry) GetValue(name string) (reflect.Value, error) {
	t, err := r.GetType(name)
	if err != nil {
		return reflect.Value{}, err
	}
	return reflect.New(t).Elem(), nil
}

// GetInterface returns an interface to a new reflect.Value of reflect.Type
// registered under specified name or an error.
func (r *Registry) GetInterface(name string) (interface{}, error) {
	t, err := r.GetType(name)
	if err != nil {
		return reflect.Value{}, err
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
	sort.Strings(names)
	return names
}
