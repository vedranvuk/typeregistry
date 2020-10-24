// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package typeregistry

import (
	"errors"
	"reflect"
	"testing"
)

type Foo struct {
	Val int
}

type PFoo *Foo

func TestRegister(t *testing.T) {

	type Test struct {
		Err      error
		Expected string
		Type     interface{}
	}

	var (
		v          = Foo{}
		pv         = &v
		ppv        = &pv
		tpv   PFoo = &v
		ptpv       = &tpv
		ntpv  PFoo = nil
		pntpv      = &ntpv
	)

	tests := []Test{
		{ErrInvalidParam, "", nil},
		{nil, "", v},
		{nil, "", pv},
		{nil, "", ppv},
		{ErrDuplicateEntry, "", tpv},
		{ErrDuplicateEntry, "", ptpv},
		{nil, "", ntpv},
		{nil, "", pntpv},
	}

	reg := New()
	for _, item := range tests {
		if err := reg.Register(item.Type); !errors.Is(err, item.Err) {
			t.Fatal(err)
		}
	}
}

func TestUnregister(t *testing.T) {

	v := &Foo{}

	reg := New()
	if err := reg.Register(v); err != nil {
		t.Fatal(err)
	}

	if err := reg.Unregister(GetLongTypeName(v)); err != nil {
		t.Fatal(err)
	}
	if err := reg.Unregister(GetLongTypeName(v)); !errors.Is(err, ErrNotFound) {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {

	v := &Foo{}

	reg := New()
	if err := reg.Register(v); err != nil {
		t.Fatal(err)
	}

	rt, err := reg.GetType(GetLongTypeName(v))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(reflect.ValueOf(v).Type(), rt) {
		t.Fatal("Invalid type returned.")
	}

	rv, err := reg.GetValue(GetLongTypeName(v))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(reflect.ValueOf(v).Type(), rv.Type()) {
		t.Fatal("Invalid value returned.")
	}

	ri, err := reg.GetInterface(GetLongTypeName(v))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(reflect.ValueOf(v).Type(), reflect.ValueOf(ri).Type()) {
		t.Fatal("Invalid interface returned.")
	}
}
