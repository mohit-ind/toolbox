package shell

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

// CommandTestSuite extends testify's Suite
type CommandTestSuite struct {
	suite.Suite
}

func (cts *CommandTestSuite) execAndCapture(f func() error) (string, error) {
	// Save os.Stdout to a var, and switch it with a buffer writer
	oldStdout := os.Stdout
	reader, writer, pipeErr := os.Pipe()
	cts.NoError(pipeErr, "OS Pipe should have been created")
	os.Stdout = writer

	// Execute the function, save the error
	err := f()

	// Read the custom writer's buffer and restore os.Stdout
	writer.Close()
	out, readErr := ioutil.ReadAll(reader)
	cts.NoError(readErr, "The buffer should be readable")
	os.Stdout = oldStdout

	// return captured output and error
	return string(out), err
}

func (cts *CommandTestSuite) tempFile(content string) string {
	tmpfile, err := ioutil.TempFile(os.TempDir(), "devops-testing")
	cts.NoError(err, "Test file should have been created")

	f, err := os.OpenFile(tmpfile.Name(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	cts.NoError(err, "Test file should have been opened for write")
	_, err = f.WriteString(content)
	cts.NoError(err, "Content should have been written into the test file")
	cts.NoError(tmpfile.Close(), "Test file should have been closed")
	cts.NoError(os.Chmod(tmpfile.Name(), 0700), "Test file should be executable now")
	return tmpfile.Name()
}

func (cts *CommandTestSuite) TestCommandWithStandardOuts() {
	out, err := cts.execAndCapture(func() error {
		cmd := NewCommandWithStandardOuts("echo", "cheese")
		return cmd.Run()
	})
	cts.NoError(err)
	cts.Contains(out, "cheese")
}

func (cts *CommandTestSuite) TestSetDir() {
	out, err := cts.execAndCapture(func() error {
		cmd := NewCommandWithStandardOuts("pwd")
		cmd.SetDir("/etc")
		return cmd.Run()
	})
	cts.NoError(err)
	cts.Contains(out, "/etc")
}

func (cts *CommandTestSuite) TestAppendEnvs() {
	out, err := cts.execAndCapture(func() error {
		cmd := NewCommandWithStandardOuts("env")
		cmd.AppendEnvs("ENV1=val1", "ENV2=val2")
		return cmd.Run()
	})
	cts.NoError(err)
	cts.Contains(out, "ENV1=val1")
	cts.Contains(out, "ENV2=val2")
	cts.Contains(out, "PATH")
}

func (cts *CommandTestSuite) TestSetStdin() {
	testFile := cts.tempFile(`
	#!/usr/bin/env sh
	cat -
	`)
	cmd := NewCommand("bash", testFile)
	cmd.SetStdin(strings.NewReader("cheese"))
	status, out, err := cmd.RunAndReturnResult()
	cts.NoError(err)
	cts.Equal(0, status)
	cts.Contains(out, "cheese")
}

func (cts *CommandTestSuite) TestRunAndReturnResult() {
	testInfoFile := cts.tempFile(`
	#!/usr/bin/env sh
	echo cheese
	`)
	cts.T().Logf("Test info file: %s", testInfoFile)
	defer os.Remove(testInfoFile)

	testErrorFile := cts.tempFile(`
	#!/usr/bin/env sh
	exit 42
	`)
	cts.T().Logf("Test error file: %s", testErrorFile)
	defer os.Remove(testErrorFile)

	tests := []struct {
		name         string
		cmd          *Command
		wantExitCode int
		wantOut      string
		wantErr      bool
	}{
		{
			name:         "invalid command",
			cmd:          NewCommand(""),
			wantExitCode: -1,
			wantErr:      true,
		},
		{
			name:         "env command",
			cmd:          NewCommand("env"),
			wantExitCode: 0,
			wantErr:      false,
		},
		{
			name:         "not existing executable",
			cmd:          NewCommand("bash", "not_existing_executable.sh"),
			wantExitCode: 127,
			wantErr:      true,
		},
		{
			name:         "output cheese",
			cmd:          NewCommand("bash", testInfoFile),
			wantExitCode: 0,
			wantOut:      "cheese",
			wantErr:      false,
		},
		{
			name:         "exit 42",
			cmd:          NewCommand("bash", testErrorFile),
			wantExitCode: 42,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		cts.T().Logf("TestRunAndReturnResult: %s", tt.name)
		gotExitCode, output, err := tt.cmd.RunAndReturnResult()
		if tt.wantErr {
			cts.Error(err)
			return
		}
		cts.Equal(tt.wantExitCode, gotExitCode)
		if tt.wantOut != "" {
			cts.Contains(output, tt.wantOut)
		}
	}
}

// TestShellCommand runs the whole test suite
func TestShellCommand(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}
