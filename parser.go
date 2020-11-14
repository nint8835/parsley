package parsley

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/shlex"
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

// RunCommand parses the content of a specific message and runs the associated command, if found.
func (parser *Parser) RunCommand(message *discordgo.MessageCreate) error {
	if !strings.HasPrefix(message.Content, parser.prefix) {
		return nil
	}

	trimmedCommand := strings.TrimPrefix(message.Content, parser.prefix)
	arguments, err := shlex.Split(trimmedCommand)
	if err != nil {
		return fmt.Errorf("error parsing arguments: %w", err)
	}

	if len(arguments) == 0 {
		return fmt.Errorf("error running command: %w", ErrNoCommandProvided)
	}

	command, ok := parser.commands[arguments[0]]
	if !ok {
		return fmt.Errorf("error running command: %w", ErrUnknownCommand)
	}

	argsParamType := reflect.TypeOf(command.handler).In(1)
	argsParamValue := reflect.New(argsParamType).Elem()

	for index := 0; index < argsParamValue.NumField(); index++ {
		field := argsParamValue.Field(index)
		var value string

		if index >= len(arguments)-1 {
			defaultVal, ok := argsParamType.Field(index).Tag.Lookup("default")
			if !ok {
				return fmt.Errorf("error parsing arguments: %w", ErrRequiredArgumentMissing)
			}
			value = defaultVal
		} else {
			value = arguments[index+1]
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int:
			intVal, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing arguments: %w", err)
			}
			field.SetInt(int64(intVal))
		case reflect.Float64:
			floatVal, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("error parsing arguments: %w", err)
			}
			field.SetFloat(floatVal)
		}
	}

	reflect.ValueOf(command.handler).Call([]reflect.Value{reflect.ValueOf(message), argsParamValue})

	return nil
}

// RegisterHandler registers a simpler handler on a discordgo session to automatically parse incoming messages for you.
func (parser *Parser) RegisterHandler(session *discordgo.Session) {
	session.AddHandler(func(message *discordgo.MessageCreate) {
		err := parser.RunCommand(message)

		if err != nil {
			_, err = session.ChannelMessageSend(
				message.ChannelID,
				fmt.Sprintf("An error occurred running your command:\n```\n%s\n```", err.Error()),
			)
			if err != nil {
				log.Fatalf("Failed to send error message: %s", err.Error())
			}
		}
	})
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
