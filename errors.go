package main

import "errors"

// ErrHandlerNotFunction occurs when a provided handler is not a function.
var ErrHandlerNotFunction error = errors.New("provided command handler is not a function")

// ErrHandlerInvalidParameterCount occurs when a provided handler does not expect two parameters.
var ErrHandlerInvalidParameterCount error = errors.New("provided command handler expects incorrect number of parameters")
