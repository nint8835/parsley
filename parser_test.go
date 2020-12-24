package parsley

import (
	"errors"
	"strconv"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-test/deep"
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

func TestRunCommandWithNotMatchingPrefix(t *testing.T) {
	parser := New("TEST")

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: "message"}})
	if err != nil {
		t.Errorf("running command returned unexpected error")
	}
}

func TestRunCommandWithSyntaxError(t *testing.T) {
	parser := New(".")

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ". \""}})
	if err == nil {
		t.Errorf("running command did not return error")
	}
}

func TestRunCommandWithUnknownCommand(t *testing.T) {
	parser := New(".")

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".unknown"}})
	if errors.Unwrap(err) != ErrUnknownCommand {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithMissingRequiredArgument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			RequiredArg string
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test"}})
	if errors.Unwrap(err) != ErrRequiredArgumentMissing {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithInvalidBoolArgument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Arg bool
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test ABC"}})
	if errors.Unwrap(errors.Unwrap(err)) != strconv.ErrSyntax {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithInvalidIntArgument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			IntArg int
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test ABC"}})
	if errors.Unwrap(errors.Unwrap(err)) != strconv.ErrSyntax {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithInvalidInt8Argument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Arg int8
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test ABC"}})
	if errors.Unwrap(errors.Unwrap(err)) != strconv.ErrSyntax {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithInvalidInt16Argument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Arg int16
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test ABC"}})
	if errors.Unwrap(errors.Unwrap(err)) != strconv.ErrSyntax {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithInvalidInt32Argument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Arg int32
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test ABC"}})
	if errors.Unwrap(errors.Unwrap(err)) != strconv.ErrSyntax {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithInvalidInt64Argument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Arg int64
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test ABC"}})
	if errors.Unwrap(errors.Unwrap(err)) != strconv.ErrSyntax {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithInvalidUintArgument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Arg uint
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test ABC"}})
	if errors.Unwrap(errors.Unwrap(err)) != strconv.ErrSyntax {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithInvalidUint8Argument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Arg uint8
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test ABC"}})
	if errors.Unwrap(errors.Unwrap(err)) != strconv.ErrSyntax {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithInvalidUint16Argument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Arg uint16
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test ABC"}})
	if errors.Unwrap(errors.Unwrap(err)) != strconv.ErrSyntax {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithInvalidUint32Argument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Arg uint32
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test ABC"}})
	if errors.Unwrap(errors.Unwrap(err)) != strconv.ErrSyntax {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithInvalidUint64Argument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Arg uint64
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test ABC"}})
	if errors.Unwrap(errors.Unwrap(err)) != strconv.ErrSyntax {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithInvalidFloat32Argument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Arg float32
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test ABC"}})
	if errors.Unwrap(errors.Unwrap(err)) != strconv.ErrSyntax {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithInvalidFloat64Argument(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			FloatArg float64
		},
	) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test ABC"}})
	if errors.Unwrap(errors.Unwrap(err)) != strconv.ErrSyntax {
		t.Errorf("running command did not return correct error")
	}
}

func TestRunCommandWithDefaultArgumentValue(t *testing.T) {
	parser := New(".")
	parser.NewCommand("test", "", func(
		message *discordgo.MessageCreate,
		args struct {
			DefaultArg string `default:"test"`
		},
	) {
		if args.DefaultArg != "test" {
			t.Errorf("handler was not passed correct value for default arg")
		}
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ".test"}})
	if err != nil {
		t.Errorf("running command returned unexpected error")
	}
}

func TestRunCommandWithEmptyCommandName(t *testing.T) {
	parser := New(".")
	parser.NewCommand("", "", func(message *discordgo.MessageCreate, args struct{}) {})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: "."}})
	if err != nil {
		t.Errorf("running command returned unexpected error")
	}
}

func TestRunCommandWithRequiredValProvidedAsKwarg(t *testing.T) {
	parser := New(".")
	parser.NewCommand("", "", func(message *discordgo.MessageCreate, args struct {
		Arg string
	}) {
		if args.Arg != "kwargval" {
			t.Errorf("handler was not passed correct value for kwarg")
		}
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ". Arg=kwargval"}})
	if err != nil {
		t.Errorf("running command returned unexpected error")
	}
}

func TestRunCommandWithKwargInMiddle(t *testing.T) {
	parser := New(".")
	parser.NewCommand("", "", func(message *discordgo.MessageCreate, args struct {
		Arg1 string
		Arg2 string
		Arg3 string
	}) {
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ". A Arg2=B C"}})
	if !errors.Is(err, ErrKwargsMustBeAtEnd) {
		t.Errorf("running command returned incorrect error")
	}
}

func TestRunCommandWithKwargAndDefault(t *testing.T) {
	parser := New(".")
	parser.NewCommand("", "", func(message *discordgo.MessageCreate, args struct {
		Arg1 string `default:"A"`
		Arg2 string `default:"B"`
	}) {
		if args.Arg2 != "overridden" {
			t.Errorf("handler was not provided correct value for kwarg")
		}
	})

	err := parser.RunCommand(&discordgo.MessageCreate{Message: &discordgo.Message{Content: ". Arg2=overridden"}})
	if err != nil {
		t.Errorf("running command returned unexpected error")
	}
}

func TestGetCommandWithUnknownCommand(t *testing.T) {
	parser := New("")

	_, err := parser.GetCommand("")
	if !errors.Is(err, ErrUnknownCommand) {
		t.Errorf("function did not return expected error")
	}
}

func TestGetCommandWithCommandWithRequiredArg(t *testing.T) {
	parser := New("")
	parser.NewCommand("", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Test string
		}) {
	})

	command, err := parser.GetCommand("")
	if err != nil {
		t.Errorf("got unexpected error")
	}

	if diff := deep.Equal(command, CommandDetails{
		Name:        "",
		Description: "",
		Arguments: []ArgumentDetails{
			{
				Name:        "Test",
				Type:        "string",
				Description: "No description provided.",
				Required:    true,
				Default:     "",
			},
		},
	}); diff != nil {
		t.Error(diff)
	}
}

func TestGetCommandWithCommandWithDefaultArg(t *testing.T) {
	parser := New("")
	parser.NewCommand("", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Test float64 `default:"1.25"`
		}) {
	})

	command, err := parser.GetCommand("")
	if err != nil {
		t.Errorf("got unexpected error")
	}

	if diff := deep.Equal(command, CommandDetails{
		Name:        "",
		Description: "",
		Arguments: []ArgumentDetails{
			{
				Name:        "Test",
				Type:        "float64",
				Description: "No description provided.",
				Required:    false,
				Default:     "1.25",
			},
		},
	}); diff != nil {
		t.Error(diff)
	}
}

func TestGetCommandWithCommandWithArgumentWithDescription(t *testing.T) {
	parser := New("")
	parser.NewCommand("", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Test string `description:"Test"`
		}) {
	})

	command, err := parser.GetCommand("")
	if err != nil {
		t.Errorf("got unexpected error")
	}

	if diff := deep.Equal(command, CommandDetails{
		Name:        "",
		Description: "",
		Arguments: []ArgumentDetails{
			{
				Name:        "Test",
				Type:        "string",
				Description: "Test",
				Required:    true,
				Default:     "",
			},
		},
	}); diff != nil {
		t.Error(diff)
	}
}

func TestGetCommandsWithMultipleCommands(t *testing.T) {
	parser := New("")
	parser.NewCommand("1", "", func(
		message *discordgo.MessageCreate,
		args struct {
		}) {
	})
	parser.NewCommand("2", "", func(
		message *discordgo.MessageCreate,
		args struct {
		}) {
	})

	commands := parser.GetCommands()

	command1, err := parser.GetCommand("1")
	if err != nil {
		t.Errorf("got unexpected error")
	}

	command2, err := parser.GetCommand("2")
	if err != nil {
		t.Errorf("got unexpected error")
	}

	if diff := deep.Equal(commands, []CommandDetails{
		command1, command2,
	}); diff != nil {
		t.Error(diff)
	}
}

func TestGetCommandsWithNoCommands(t *testing.T) {
	parser := New("")

	commands := parser.GetCommands()

	if diff := deep.Equal(commands, []CommandDetails{}); diff != nil {
		t.Error(diff)
	}
}
