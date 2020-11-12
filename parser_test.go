package main

import "testing"

func TestNewParser(t *testing.T) {
	parser := New("test")
	if parser.prefix != "test" {
		t.Errorf("parser was returned with incorrect prefix %s", parser.prefix)
	}
	if len(parser.commands) != 0 {
		t.Error("parser was returned with non-empty initial command list")
	}
}
