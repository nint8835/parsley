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

func TestGetCommandsWithCommandWithRequiredArg(t *testing.T) {
	parser := New("")
	parser.NewCommand("", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Test string
		}) {
	})

	commands := parser.GetCommands()

	if diff := deep.Equal(commands, []CommandDetails{{
		Name:        "",
		Description: "",
		Arguments: []ArgumentDetails{
			{
				Name:     "Test",
				Type:     "string",
				Required: true,
				Default:  "",
			},
		},
	}}); diff != nil {
		t.Error(diff)
	}
}

func TestGetCommandsWithCommandWithDefaultArg(t *testing.T) {
	parser := New("")
	parser.NewCommand("", "", func(
		message *discordgo.MessageCreate,
		args struct {
			Test float64 `default:"1.25"`
		}) {
	})

	commands := parser.GetCommands()

	if diff := deep.Equal(commands, []CommandDetails{{
		Name:        "",
		Description: "",
		Arguments: []ArgumentDetails{
			{
				Name:     "Test",
				Type:     "float64",
				Required: false,
				Default:  "1.25",
			},
		},
	}}); diff != nil {
		t.Error(diff)
	}
}

func TestGetCommandsWithCommandMultipleCommands(t *testing.T) {
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

	if diff := deep.Equal(commands, []CommandDetails{
		{
			Name:        "1",
			Description: "",
			Arguments:   []ArgumentDetails{},
		},
		{
			Name:        "2",
			Description: "",
			Arguments:   []ArgumentDetails{},
		},
	}); diff != nil {
		t.Error(diff)
	}
}
