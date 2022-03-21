package cli

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/suite"
)

///////////
// Suite //
///////////

// CLITestSuite extends testify's Suite.
type CLITestSuite struct {
	suite.Suite
}

func (cts *CLITestSuite) getTestCommand(name string) *Command {
	cmd := NewCommand(name)

	cts.NotNil(cmd, "Command should have been created")
	cts.Equalf(name, cmd.Name, "The new Command's name should be %s", name)
	cts.Nil(cmd.Execute([]string{}), "The new Command should have an empty task which returns nil")
	return cmd
}

func (cts *CLITestSuite) TestNewCommand() {
	cts.getTestCommand("test-command")
}

func (cts *CLITestSuite) TestWithAliases() {
	cmd := cts.getTestCommand("test-command").WithAliases("aliasA", "aliasB")
	cts.Contains(cmd.Aliases, "aliasA", "aliasA should have been added to the new Command as an alias")
	cts.Contains(cmd.Aliases, "aliasB", "aliasB should have been added to the new Command as an alias")
}

func (cts *CLITestSuite) TestWithTask() {
	testList := []string{}

	testTask := func(args []string) error {
		testList = args
		return errors.New("Test Error 1")
	}

	cmd := cts.getTestCommand("test-command").WithTask(testTask)

	err := cmd.Execute([]string{"arg1", "arg2"})
	cts.Error(err, "Execute should have returned a test error")
	cts.Contains(err.Error(), "Test Error 1", "The error returned by Execute should contain Test Error 1")
	cts.Contains(testList, "arg1", "The test list should contain arg1 after the Command is executed")
	cts.Contains(testList, "arg2", "The test list should contain arg2 after the Command is executed")
}

func (cts *CLITestSuite) TestWithSubCommand() {
	testList := []string{}

	testTask1 := func(args []string) error {
		testList = args
		return errors.New("Test Error 1")
	}

	testTask2 := func(args []string) error {
		testList = args
		return errors.New("Test Error 2")
	}

	testTask3 := func(args []string) error {
		testList = args
		return errors.New("Test Error 3")
	}

	subCmd1 := cts.getTestCommand("sub-command1").WithTask(testTask1)
	subCmd2 := cts.getTestCommand("sub-command2").WithTask(testTask2).WithAliases("sub2")
	subCmd3 := cts.getTestCommand("sub-command3").WithTask(testTask3).WithSubCommands(subCmd2)
	cmd := cts.getTestCommand("test-command").WithSubCommands(subCmd1, subCmd3)

	// sub-command2 cannot be called directly on root level, so the root command's task should be executed
	cts.Nil(cmd.Execute([]string{"sub-command2"}), "The root command should not return any error")
	cts.Empty(testList, "The test list should be empty now")

	err := cmd.Execute([]string{"sub-command1", "arg2"})
	cts.Error(err, "Test sub-command1 should have returned an error")
	cts.Contains(err.Error(), "Test Error 1", "The error returned from sub-command1 should contain Test Error 1")
	cts.Contains(testList, "arg2", "The test list should contain arg2 after sub-command1 is executed")
	cts.NotContains(testList, "sub-command1", "The test list shouldn't contain the sub-command1 name")

	err = cmd.Execute([]string{"sub-command3", "arg3"})
	cts.Error(err, "Test sub-command3 should have returned an error")
	cts.Contains(err.Error(), "Test Error 3", "The error returned from sub-command3 should contain Test Error 3")
	cts.Contains(testList, "arg3", "The test list should contain arg3 after sub-command3 is executed")
	cts.NotContains(testList, "sub-command3", "The test list shouldn't contain the sub-command3 name")
	cts.NotContains(testList, "arg2", "The test list shouldn't contain arg2 after sub-command3 executed")

	err = cmd.Execute([]string{"sub-command3", "sub2", "arg4"})
	cts.Error(err, "Test sub-command3 should have returned an error")
	cts.Contains(err.Error(), "Test Error 2", "The error returned from sub-command3 should contain Test Error 2")
	cts.Contains(testList, "arg4", "The test list should contain arg4 after sub-command3 is executed")
	cts.NotContains(testList, "sub-command3", "The test list shouldn't contain the sub-command3 name")
	cts.NotContains(testList, "arg2", "The test list shouldn't contain arg2 after sub-command3 executed")
	cts.NotContains(testList, "arg3", "The test list shouldn't contain arg3 after sub-command3 executed")
}

func (cts *CLITestSuite) TestEndWithMessage_NoCommand() {
	// Switch os.Stdout with a custom writer
	oldStdout := os.Stdout
	reader, writer, pipeErr := os.Pipe()
	cts.NoError(pipeErr, "OS Pipe should have been created")
	os.Stdout = writer

	// Create an end task and call it
	f := EndWithMessage("Error 1")
	cts.NotNil(f, "A Task should have been created")
	err := f([]string{})
	cts.Error(err, "An error should have been generated by EndWithMessage")
	cts.Contains(err.Error(), "Command is missing")

	// Read the custom writer's buffer and restore os.Stdout
	writer.Close()
	out, readErr := ioutil.ReadAll(reader)
	cts.NoError(readErr, "The buffer should have been red")
	os.Stdout = oldStdout
	cts.Contains(string(out), "Error 1", "The output should contain Error 1")
}

func (cts *CLITestSuite) TestEndWithMessage_InvalidCommand() {
	// Switch os.Stdout with a custom writer
	oldStdout := os.Stdout
	reader, writer, pipeErr := os.Pipe()
	cts.NoError(pipeErr, "OS Pipe should have been created")
	os.Stdout = writer

	// Create an end task and call it
	f := EndWithMessage("Error 2")
	cts.NotNil(f, "A Task should have been created")
	err := f([]string{"AnInvalidCommand"})
	cts.Error(err, "An error should have been generated by EndWithMessage")
	cts.Contains(err.Error(), "Invalid Command: AnInvalidCommand")

	// Read the custom writer's buffer and restore os.Stdout
	writer.Close()
	out, readErr := ioutil.ReadAll(reader)
	cts.NoError(readErr, "The buffer should have been red")
	os.Stdout = oldStdout
	cts.Contains(string(out), "Error 2", "The output should contain Error 1")
}

// TestCLI runs the whole test suite
func TestCLI(t *testing.T) {
	suite.Run(t, new(CLITestSuite))
}
