package shell

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

// Command represents a shell command
type Command struct {
	cmd *exec.Cmd
}

// NewCommand creates a new command
func NewCommand(name string, args ...string) *Command {
	return &Command{
		cmd: exec.Command(name, args...),
	}
}

// NewWCommandithStandardOuts - same as NewCommand, but sets the command's
// stdout and stderr to the standard (OS) out (os.Stdout) and err (os.Stderr).
func NewCommandWithStandardOuts(name string, args ...string) *Command {
	return NewCommand(name, args...).SetStdout(os.Stdout).SetStderr(os.Stderr)
}

// SetDir sets the execution directory of the command.
func (c *Command) SetDir(dir string) *Command {
	c.cmd.Dir = dir
	return c
}

// SetEnvs environment variables for the command in the "key=value" format.
func (c *Command) SetEnvs(envs ...string) *Command {
	c.cmd.Env = envs
	return c
}

// AppendEnvs - appends the envs to the current os.Environ()
// Calling this multiple times will NOT appens the envs one by one,
// only the last "envs" set will be appended to os.Environ()!
func (c *Command) AppendEnvs(envs ...string) *Command {
	return c.SetEnvs(append(os.Environ(), envs...)...)
}

// SetStdin sets the standard input for the command.
func (c *Command) SetStdin(in io.Reader) *Command {
	c.cmd.Stdin = in
	return c
}

// SetStdout sets the standard output of the command.
func (c *Command) SetStdout(out io.Writer) *Command {
	c.cmd.Stdout = out
	return c
}

// SetStderr sets the standard error of the command.
func (c *Command) SetStderr(err io.Writer) *Command {
	c.cmd.Stderr = err
	return c
}

// Run executes the command.
func (c Command) Run() error {
	return c.cmd.Run()
}

// RunAndReturnResult executes the command and returns its exit code, combined stdout and stderr and an optional error.
func (c Command) RunAndReturnResult() (int, string, error) {
	outBytes, err := c.cmd.CombinedOutput()
	return c.cmd.ProcessState.ExitCode(), strings.TrimSpace(string(outBytes)), err
}
