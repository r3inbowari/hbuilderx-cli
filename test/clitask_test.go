package test

import (
	"meiwobuxing"
	"testing"
)

func TestFt(t *testing.T) {
	a := meiwobuxing.NewTask()
	a.AddProcess("12", func() error {
		println("1")
		return nil
	}).AddProcess("23", func() error {
		println("2")
		return nil
	})

	a.Start()
}
