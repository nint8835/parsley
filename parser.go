package main

import "reflect"

// Command represents an individual Discord command.
type Command struct {
	description string
	handler     interface{}
}

// Parser represents a parser for Discord commands.
type Parser struct {
	prefix   string
	commands map[string]Command
}

// New creates a new Parsley parser.
func New(prefix string) *Parser {
	return &Parser{
		prefix, make(map[string]Command, 0),
	}
}

func _ValidateHandler(handler interface{}) error {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		return ErrHandlerNotFunction
	}
	if handlerType.NumIn() != 2 {
		return ErrHandlerInvalidParameterCount
	}
	return nil
}

// NewCommand registers a new command with the command parser.
func (parser *Parser) NewCommand(name, description string, handler interface{}) {
	command := Command{description, handler}
	parser.commands[name] = command
}
