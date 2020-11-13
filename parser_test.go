package main

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestNewParser(t *testing.T) {
	parser := New("test")
	if parser.prefix != "test" {
		t.Errorf("parser was returned with incorrect prefix %s", parser.prefix)
	}
	if len(parser.commands) != 0 {
		t.Error("parser was returned with non-empty initial command list")
	}
}

func TestNewCommandWithInvalidHandlerType(t *testing.T) {
	parser := New("")

	err := parser.NewCommand("", "", "abc123")
	if errors.Unwrap(err) != ErrHandlerNotFunction {
		t.Error("parser did not return correct error")
	}
}

func TestNewCommandWithInvalidParameterCount(t *testing.T) {
	parser := New("")

	err := parser.NewCommand("", "", func() {})
	if errors.Unwrap(err) != ErrHandlerInvalidParameterCount {
		t.Error("parser did not return correct error")
	}
}

func TestNewCommandWithInvalidFirstParameterType(t *testing.T) {
	parser := New("")

	err := parser.NewCommand("", "", func(a *string, b struct{}) {})
	if errors.Unwrap(err) != ErrHandlerInvalidFirstParameterType {
		t.Error("parser did not return correct error")
	}
}

func TestNewCommandWithFirstParameterNotPointer(t *testing.T) {
	parser := New("")

	err := parser.NewCommand("", "", func(a discordgo.MessageCreate, b struct{}) {})
	if errors.Unwrap(err) != ErrHandlerInvalidFirstParameterType {
		t.Error("parser did not return correct error")
	}
}

func TestNewCommandWithInvalidSecondParameterType(t *testing.T) {
	parser := New("")

	err := parser.NewCommand("", "", func(a *discordgo.MessageCreate, b string) {})
	if errors.Unwrap(err) != ErrHandlerInvalidSecondParameterType {
		t.Error("parser did not return correct error")
	}
}

func TestNewCommandWithValidData(t *testing.T) {
	parser := New("")

	err := parser.NewCommand("", "", func(a *discordgo.MessageCreate, b struct{}) {})
	if err != nil {
		t.Errorf("adding command returned unexpected error")
	}
}
