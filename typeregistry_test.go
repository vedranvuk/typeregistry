// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package typeregistry

import (
	"errors"
	"reflect"
	"testing"
)

func TestTypeRegistry(t *testing.T) {
	type Foo struct{ _ bool }
	type Bar struct{ _ bool }
	reg := New()
	foo := &Foo{}
	bar := &Bar{}
	fooname := GetLongTypeName(foo)
	barname := GetLongTypeName(bar)

	if err := reg.Register(foo); err != nil {
		t.Fatal(err)
	}
	if err := reg.RegisterNamed(barname, bar); err != nil {
		t.Fatal(err)
	}
	if err := reg.Register(bar); !errors.Is(err, ErrDuplicateEntry) {
		t.Fatal("Register failed.")
	}
	if err := reg.RegisterNamed("", foo); err != ErrInvalidParam {
		t.Fatal(err)
	}
	if err := reg.RegisterNamed(barname, nil); err != ErrInvalidParam {
		t.Fatal(err)
	}

	rn := reg.RegisteredNames()
	exprn := []string{
		barname,
		fooname,
	}
	for idx, name := range rn {
		if exprn[idx] != name {
			t.Fatal("RegisteredNames failed.")
		}
	}

	rt, err := reg.GetType(GetLongTypeName(foo))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = reg.GetType("baz"); !errors.Is(err, ErrNotFound) {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(reflect.ValueOf(foo).Type(), rt) {
		t.Fatal("Invalid type returned.")
	}

	rv, err := reg.GetValue(GetLongTypeName(foo))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = reg.GetValue("baz"); !errors.Is(err, ErrNotFound) {
		t.Fatal("GetValue failed.")
	}
	if !reflect.DeepEqual(reflect.ValueOf(foo).Type(), rv.Type()) {
		t.Fatal("Invalid value returned.")
	}

	ri, err := reg.GetInterface(GetLongTypeName(foo))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = reg.GetInterface("baz"); !errors.Is(err, ErrNotFound) {
		t.Fatal("GetInterface failed.")
	}
	if !reflect.DeepEqual(reflect.ValueOf(foo).Type(), reflect.ValueOf(ri).Type()) {
		t.Fatal("Invalid interface returned.")
	}

	if err := reg.Unregister("foo"); !errors.Is(err, ErrNotFound) {
		t.Fatal("Unregister failed.")
	}
	if err := reg.Unregister(fooname); err != nil {
		t.Fatal()
	}
	if err := reg.Unregister(barname); err != nil {
		t.Fatal()
	}
	if len(reg.RegisteredNames()) != 0 {
		t.Fatal("Unregister failed.")
	}
}

func TestGetLongTypeName(t *testing.T) {
	type alias int
	type palias *alias
	var valias alias = 42
	var vpalias palias = &valias
	if GetLongTypeName(nil) != "" {
		t.Fatal("GetLongTypeName failed.")
	}
	if GetLongTypeName(int(0)) != "int" {
		t.Fatal("GetLongTypeName failed.")
	}
	if GetLongTypeName(string("")) != "string" {
		t.Fatal("GetLongTypeName failed.")
	}
	if GetLongTypeName(valias) != "github.com/vedranvuk/typeregistry/typeregistry.alias" {
		t.Fatal("GetLongTypeName failed.")
	}
	if GetLongTypeName(alias(0)) != "github.com/vedranvuk/typeregistry/typeregistry.alias" {
		t.Fatal("GetLongTypeName failed.")
	}
	if GetLongTypeName(alias(42)) != "github.com/vedranvuk/typeregistry/typeregistry.alias" {
		t.Fatal("GetLongTypeName failed.")
	}
	if GetLongTypeName(vpalias) != "*github.com/vedranvuk/typeregistry/typeregistry.alias" {
		t.Fatal("GetLongTypeName failed.")
	}
	if GetLongTypeName(palias(nil)) != "github.com/vedranvuk/typeregistry/typeregistry.palias" {
		t.Fatal("GetLongTypeName failed.")
	}
	if GetLongTypeName(palias(vpalias)) != "*github.com/vedranvuk/typeregistry/typeregistry.alias" {
		t.Fatal("GetLongTypeName failed.")
	}
}
