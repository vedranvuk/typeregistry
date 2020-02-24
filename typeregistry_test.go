// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package typeregistry

import (
	"testing"
	"time"

	"github.com/vedranvuk/testex"
)

type TestA = struct {
	Name string
}

type TestB = struct {
	Age int
}

type TestC = struct {
	when time.Time
}

func TestRegistry(t *testing.T) {
	r := New()
	if err := r.Register("TestA", TestA{}); err != nil {
		t.Fatal(err)
	}
	if err := r.Register("TestB", TestB{}); err != nil {
		t.Fatal(err)
	}
	if err := r.Register("TestC", TestC{}); err != nil {
		t.Fatal(err)
	}

	v, err := r.Get("TestC")
	if err != nil {
		t.Fatal(err)
	}
	rv, ok := v.(TestC)
	if !ok {
		t.Fatal("fail")
	}

	if testex.Verbose {
		t.Log(rv)
	}
}
