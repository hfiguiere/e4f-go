package main

import "testing"

func TestE4f(t *testing.T) {
	e4fDb := Parse("samples/export-Roll-20130630_203650.xml")

	e4fDb.buildMaps()

	if l := len(e4fDb.Cameras); l != 1 {
		t.Errorf("Found %d cameras, expected 1", l)
	}
	cam := e4fDb.Cameras[0]
	if id := cam.MakeId; id != 13 {
		t.Errorf("Camera MakeId expect is 13. Found %d", id)
	}
	if l := len(e4fDb.Makes); l != 2 {
		t.Errorf("Found %d makes, expected 2", l)
	}

	if l := len(e4fDb.Exposures); l != 37 {
		t.Errorf("Found %d exposures", l)
	}
}
