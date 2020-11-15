package parsley

import (
	"fmt"
	"log"
	"reflect"
	"sort"
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

// ArgumentDetails represents the details of an individual command argument.
type ArgumentDetails struct {
	Name     string
	Type     string
	Required bool
	Default  string
}

// CommandDetails represents the parsed details of an individual command.
type CommandDetails struct {
	Name        string
	Description string
	Arguments   []ArgumentDetails
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

	arguments, err := shlex.Split(message.Content)
	arguments[0] = strings.TrimPrefix(arguments[0], parser.prefix)
	if err != nil {
		return fmt.Errorf("error parsing arguments: %w", err)
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
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("error parsing arguments: %w", err)
			}
			field.SetBool(boolVal)
		case reflect.Int, reflect.Int64:
			intVal, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error parsing arguments: %w", err)
			}
			field.SetInt(int64(intVal))
		case reflect.Int8:
			intVal, err := strconv.ParseInt(value, 10, 8)
			if err != nil {
				return fmt.Errorf("error parsing arguments: %w", err)
			}
			field.SetInt(intVal)
		case reflect.Int16:
			intVal, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				return fmt.Errorf("error parsing arguments: %w", err)
			}
			field.SetInt(intVal)
		case reflect.Int32:
			intVal, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return fmt.Errorf("error parsing arguments: %w", err)
			}
			field.SetInt(intVal)
		case reflect.Uint, reflect.Uint64:
			intVal, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing arguments: %w", err)
			}
			field.SetUint(intVal)
		case reflect.Uint8:
			intVal, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return fmt.Errorf("error parsing arguments: %w", err)
			}
			field.SetUint(intVal)
		case reflect.Uint16:
			intVal, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return fmt.Errorf("error parsing arguments: %w", err)
			}
			field.SetUint(intVal)
		case reflect.Uint32:
			intVal, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return fmt.Errorf("error parsing arguments: %w", err)
			}
			field.SetUint(intVal)
		case reflect.Float32:
			floatVal, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return fmt.Errorf("error parsing arguments: %w", err)
			}
			field.SetFloat(floatVal)
		case reflect.Float64:
			floatVal, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("error parsing arguments: %w", err)
			}
			field.SetFloat(floatVal)
		case reflect.String:
			field.SetString(value)
		}
	}

	reflect.ValueOf(command.handler).Call([]reflect.Value{reflect.ValueOf(message), argsParamValue})

	return nil
}

// RegisterHandler registers a simpler handler on a discordgo session to automatically parse incoming messages for you.
func (parser *Parser) RegisterHandler(session *discordgo.Session) {
	session.AddHandler(func(session *discordgo.Session, message *discordgo.MessageCreate) {
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

// GetCommands parses all registered commands and returns details related to each of them.
func (parser *Parser) GetCommands() []CommandDetails {
	commandDetails := make([]CommandDetails, 0)
	commands := make([]string, 0)
	for command := range parser.commands {
		commands = append(commands, command)
	}
	sort.Strings(commands)

	for _, commandName := range commands {
		commandObj := parser.commands[commandName]
		commandDetailsObj := CommandDetails{
			Name:        commandName,
			Description: commandObj.description,
			Arguments:   make([]ArgumentDetails, 0),
		}

		argsType := reflect.TypeOf(commandObj.handler).In(1)

		for index := 0; index < argsType.NumField(); index++ {
			arg := argsType.Field(index)

			defaultVal, hasDefault := arg.Tag.Lookup("default")
			commandDetailsObj.Arguments = append(commandDetailsObj.Arguments, ArgumentDetails{
				Name:     arg.Name,
				Type:     arg.Type.Name(),
				Required: !hasDefault,
				Default:  defaultVal,
			})
		}

		commandDetails = append(commandDetails, commandDetailsObj)
	}

	return commandDetails
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
