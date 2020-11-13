package main

import (
	"fmt"
	"reflect"

	"github.com/bwmarrin/discordgo"
)

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
	firstParam := handlerType.In(0)
	if firstParam.Kind() != reflect.Ptr || firstParam.Elem() != reflect.TypeOf(discordgo.MessageCreate{}) {
		return ErrHandlerInvalidFirstParameterType
	}
	if handlerType.In(1).Kind() != reflect.Struct {
		return ErrHandlerInvalidSecondParameterType
	}
	return nil
}

// NewCommand registers a new command with the command parser.
func (parser *Parser) NewCommand(name, description string, handler interface{}) error {
	err := _ValidateHandler(handler)
	if err != nil {
		return fmt.Errorf("invalid command handler: %w", err)
	}
	command := Command{description, handler}
	parser.commands[name] = command

	return nil
}
