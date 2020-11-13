package main

import "errors"

// ErrHandlerNotFunction occurs when a provided handler is not a function.
var ErrHandlerNotFunction error = errors.New("provided command handler is not a function")

// ErrHandlerInvalidParameterCount occurs when a provided handler does not expect two parameters.
var ErrHandlerInvalidParameterCount error = errors.New("provided command handler expects incorrect number of parameters")

// ErrHandlerInvalidFirstParameterType occurs when a provided handler does not expect a first parameter of the correct type.
var ErrHandlerInvalidFirstParameterType error = errors.New("incorrect first parameter type for handler, first parameter must be of type *discordgo.MessageCreate")

// ErrHandlerInvalidSecondParameterType occurs when a provided handler does not expect a second parameter of the correct type.
var ErrHandlerInvalidSecondParameterType error = errors.New("incorrect second parameter type for handler, second parameter must be of type struct")
