package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"

	cli "github.com/toolboxcli"
)

const VERSION = "v1.3.2"

var items = []string{
	"Hat",
	"Parrot",
	"Space Station",
	"Siege Tower",
}

func main() {
	versionCommand := cli.NewCommand("version").
		WithAliases("ver", "-v", "--version").
		WithTask(func(args []string) error {
			fmt.Println(VERSION)
			return nil
		})

	listItemsCommand := cli.NewCommand("list").
		WithAliases("l").
		WithTask(listItems)

	getItemsCommand := cli.NewCommand("get").
		WithAliases("g").
		WithTask(getItem)

	addItemsCommand := cli.NewCommand("add").
		WithAliases("a").
		WithTask(addItems)

	itemsCommand := cli.NewCommand("items").
		WithTask(cli.EndWithMessage(itemsUsage())).
		WithSubCommands(listItemsCommand, getItemsCommand, addItemsCommand)

	rootCommand := cli.NewCommand("root").
		WithSubCommands(versionCommand, itemsCommand).
		WithTask(cli.EndWithMessage(usage()))

	if err := rootCommand.Execute(os.Args[1:]); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func usage() string {
	return `
Simple CLI Example

Usage: $ app <subcommand> <arg>

Example: $ app items list

Available subcommands:
	version  - Prints the application's version
	items    - Manage Item Store
	`
}

func itemsUsage() string {
	return `
Item Store

Usage: $ app items <subcommand> <args>

Example: $ app items list

Available subcommands:
	list     - List Item Store
	get      - Get an element from the Item Store by index
	add      - Add new element(s) to the Item Store
	`
}

func listItems(args []string) error {
	fmt.Println("Item Store:")
	for index, item := range items {
		fmt.Printf("[%d][%s]\n", index, item)
	}
	return nil
}

func getItem(args []string) error {
	if len(args) < 1 {
		return errors.New("No item index provided")
	}
	val, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.Wrap(err, "Failed to get item by index")
	}
	if val > len(items)-1 {
		return errors.Errorf("Index: %d is out of range: %d", val, len(items))
	}
	fmt.Printf("[%d][%s]\n", val, items[val])
	return nil
}

func addItems(args []string) error {
	if len(args) == 0 {
		return errors.New("No item to add")
	}
	origLen := len(items)
	items = append(items, args...)
	if len(args) > 0 {
		fmt.Println("Item(s) Added:")
		for index, item := range args {
			fmt.Printf("[%d][%s]\n", origLen+index, item)
		}
	}
	return nil
}
