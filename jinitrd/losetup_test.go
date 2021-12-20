package main

import (
	"testing"

	"gopkg.in/freddierice/go-losetup.v1"
)

func TestLosetup(t *testing.T) {
	t.Log("TestLosetup")
	device, err := losetup.Attach("test.img", 0, true)
	if err != nil {
		t.Error(err)
	}
	t.Log(device.Path())
}
