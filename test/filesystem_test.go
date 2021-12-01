package test

import (
	"meiwobuxing"
	"testing"
)

func TestCal(t *testing.T) {
	meiwobuxing.InitFileSystem("res", 100)
	println(meiwobuxing.Calc("f0849618-4ab5-4b8c-bdfb-96ea85ffaa7c"))

	println(meiwobuxing.GetPath("1232"))
}
