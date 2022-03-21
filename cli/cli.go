// Package cli provides a command line argument parser and task executer
// which you can use to build simple command line interfaces for your application
package cli

import (
	"fmt"

	"github.com/pkg/errors"
)

// Command represents a terminal command
type Command struct {
	Name        string
	Aliases     []string
	Task        func(args []string) error
	SubCommands []*Command
}

// NewCommand creates a new Command with the given name
func NewCommand(name string) *Command {
	return &Command{
		Name: name,
		Task: func(args []string) error {
			return nil
		},
	}
}

// WithAliases adds the supplied strings to the Command as Aliases
func (c *Command) WithAliases(aliases ...string) *Command {
	c.Aliases = append(c.Aliases, aliases...)
	return c
}

// WithTask adds the supplied task to the Command
func (c *Command) WithTask(task func(args []string) error) *Command {
	c.Task = task
	return c
}

// WithSubCommands adds the supplied list of Commands to the Command as SubCommands
func (c *Command) WithSubCommands(subCommands ...*Command) *Command {
	c.SubCommands = append(c.SubCommands, subCommands...)
	return c
}

// Execute starts the recursive execution of the CLI
func (c *Command) Execute(args []string) error {
	// If there is 1 or more args, check if the first arg matches with any subcommand
	if len(args) > 0 {
		for _, subCommand := range c.SubCommands {
			// Check for both the Command's Name and the Aliases too
			if subCommand.Name == args[0] || match(args[0], subCommand.Aliases) {
				return subCommand.Execute(args[1:])
			}
		}
	}
	// If no match were found just execute the Task
	return c.Task(args)
}

// match check if a string is in a slice of strings
func match(s string, clues []string) bool {
	for _, clue := range clues {
		if s == clue {
			return true
		}
	}
	return false
}

// EndWithMessage returns a Task that prints out the supplied message and returns with an error
// depending on the amount of remaining args (no command / invalid command)
func EndWithMessage(msg string) func(args []string) error {
	return func(args []string) error {
		fmt.Println(msg)
		if len(args) > 0 {
			return errors.Errorf("Invalid Command: %s", args[0])
		}
		return errors.New("Command is missing")
	}
}
