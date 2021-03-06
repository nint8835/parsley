package parsley

import "errors"

// ErrHandlerNotFunction occurs when a provided handler is not a function.
var ErrHandlerNotFunction error = errors.New("provided command handler is not a function")

// ErrHandlerInvalidParameterCount occurs when a provided handler does not expect two parameters.
var ErrHandlerInvalidParameterCount error = errors.New("provided command handler expects incorrect number of parameters")

// ErrHandlerInvalidFirstParameterType occurs when a provided handler does not expect a first parameter of the correct type.
var ErrHandlerInvalidFirstParameterType error = errors.New("incorrect first parameter type for handler, first parameter must be of type *discordgo.MessageCreate")

// ErrHandlerInvalidSecondParameterType occurs when a provided handler does not expect a second parameter of the correct type.
var ErrHandlerInvalidSecondParameterType error = errors.New("incorrect second parameter type for handler, second parameter must be of type struct")

// ErrUnknownCommand occurs when the provided message or function call contains an unknown command.
var ErrUnknownCommand error = errors.New("unknown command")

// ErrRequiredArgumentMissing occurs when the provided message does not have values for all required arguments.
var ErrRequiredArgumentMissing error = errors.New("one or more required arguments were not provided")

// ErrKwargsMustBeAtEnd occurs when a user provides keyword arguments in the middle of positional arguments
var ErrKwargsMustBeAtEnd error = errors.New("keyword arguments must be provided as the last arguments")
