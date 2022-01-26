package test

import (
	"github.com/dgraph-io/badger/v3"
	"meiwobuxing"
	"testing"
)

type Hello struct {
	A string
}

func TestSGDB(t *testing.T) {
	dbOpts := badger.DefaultOptions("../db")
	err := meiwobuxing.InitDB(dbOpts)
	if err != nil {
		println(err.Error())
		return
	}

	h := Hello{A: "hello"}
	err = meiwobuxing.SetJson("go", "hello", &h)
	if err != nil {
		println(err.Error())
		t.Fail()
		return
	}

	h = Hello{A: "hello2"}
	err = meiwobuxing.SetJson("go", "hello2", &h)
	if err != nil {
		println(err.Error())
		t.Fail()
		return
	}

	h = Hello{A: "hello2"}
	err = meiwobuxing.SetJson("patch", "hello2", &h)
	if err != nil {
		println(err.Error())
		t.Fail()
		return
	}

	b := Hello{}
	err = meiwobuxing.GetJson("go", "hello", &b)
	if err != nil {
		println(err.Error())
		return
	}

	m, err := meiwobuxing.Iter("go")
	if err != nil {
		return
	}
	println(len(m))

	println(b.A)

}
