// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package typeregistry

import (
	"errors"
	"reflect"
	"testing"
)

type TestA struct {
	Name string
}

type TestB struct {
	Age int
}

type TestC struct {
	Married bool
}

type TestD struct {
	TestA
	TestB
	TestC
}

func TestRegister(t *testing.T) {

	val := &TestD{}

	reg := New()

	if err := reg.RegisterNamed("test", val); err != nil {
		t.Fatal(err)
	}

	if err := reg.Unregister("test"); err != nil {
		t.Fatal(err)
	}

	if err := reg.Register(val); err != nil {
		t.Fatal(err)
	}

	if err := reg.Register(val); !errors.Is(err, ErrDuplicateEntry) {
		t.Fatal("TestRegister failed")
	}

	if err := reg.RegisterNamed("test", val); err != nil {
		t.Fatal(err)
	}

	if err := reg.RegisterNamed("test", val); !errors.Is(err, ErrDuplicateEntry) {
		t.Fatal("TestRegister failed")
	}

	if err := reg.Unregister("INVALID"); !errors.Is(err, ErrNotFound) {
		t.Fatal(err)
	}

}

func TestGetType(t *testing.T) {

	val := &TestD{}
	rval := reflect.ValueOf(val)
	vtyp := rval.Type()

	reg := New()

	if err := reg.Register(val); err != nil {
		t.Fatal(err)
	}

	if _, err := reg.GetType("INVALID"); !errors.Is(err, ErrNotFound) {
		t.Fatal("GetType failed")
	}

	typ, err := reg.GetType(vtyp.String())
	if err != nil {
		t.Fatal(err)
	}

	if typ != vtyp {
		t.Fatal("GetType failed")
	}
}

func TestGetValue(t *testing.T) {

	val := &TestD{}
	valtyp := reflect.TypeOf(val)

	reg := New()

	if err := reg.RegisterNamed("test", val); err != nil {
		t.Fatal(err)
	}

	if _, err := reg.GetValue("INVALID"); !errors.Is(err, ErrNotFound) {
		t.Fatal("GetValue failed")
	}

	nval, err := reg.GetValue("test")
	if err != nil {
		t.Fatal(err)
	}

	if nval.Type() != valtyp {
		t.Fatal("GetValue failed")
	}
}

func TestGetInterface(t *testing.T) {

	val := &TestD{}

	reg := New()
	if err := reg.RegisterNamed("test", val); err != nil {
		t.Fatal(err)
	}

	if _, err := reg.GetInterface("INVALID"); !errors.Is(err, ErrNotFound) {
		t.Fatal("GetInterface failed")
	}

	iface, err := reg.GetInterface("test")
	if err != nil {
		t.Fatal(err)
	}

	_, ok := iface.(*TestD)
	if !ok {
		t.Fatal("GetInterface failed")
	}
}

func TestRegisteredNames(t *testing.T) {

	expected := []string{
		"*typeregistry.TestA",
		"*typeregistry.TestB",
		"*typeregistry.TestC",
		"*typeregistry.TestD",
		"typeregistry.TestA",
		"typeregistry.TestB",
		"typeregistry.TestC",
		"typeregistry.TestD",
	}

	r := New()

	if err := r.Register(TestA{}); err != nil {
		t.Fatal(err)
	}

	if err := r.Register(&TestA{}); err != nil {
		t.Fatal(err)
	}

	if err := r.Register(TestB{}); err != nil {
		t.Fatal(err)
	}

	if err := r.Register(&TestB{}); err != nil {
		t.Fatal(err)
	}

	if err := r.Register(TestC{}); err != nil {
		t.Fatal(err)
	}

	if err := r.Register(&TestC{}); err != nil {
		t.Fatal(err)
	}

	if err := r.Register(TestD{}); err != nil {
		t.Fatal(err)
	}

	if err := r.Register(&TestD{}); err != nil {
		t.Fatal(err)
	}

	names := r.RegisteredNames()
	for i := 0; i < len(names); i++ {
		if names[i] != expected[i] {
			t.Fatal("RegisteredNames failed")
		}
	}
}
