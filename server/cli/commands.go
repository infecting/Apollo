package cli

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"strings"
)

// The command type
type command struct {
	Name     string                // The Name of the command
	Aliases  []string              // List of aliases the command can be referred to by
	MinArgs  int                   // The number of args required for this command. Note that 0 args means just the command
	Help     string                // The help string for the command which will be printed if there is an error or the help command is used
	Usage	 string				   // The usage for the command which will be printed out if the command is invalid or the help command is used
	Function func(Cli, []string) string // The function which will be ran. This function takes client refrence and the args and returns the output
}

var commandList = []command{
	{
		Name:     "clear",
		Aliases:  []string{"c", "cl"},
		MinArgs:  0,
		Help:     "Clears the screen",
		Function: func(c Cli, args []string) string {
			c.Clear()
			return ""
		},
	},
	{
		Name:     "exit",
		Aliases:  []string{"quit"},
		MinArgs:  0,
		Help:     "Exists the cli",
		Function: func(c Cli, args []string) string {
			c.remove()
			return ""
		},
	},
}

/*
 * Gets the description of the command used for the help command
 * Takes an aurora object as an argument for color support
 */
func (c command) description(au aurora.Aurora) string {
	aliasesString := ""
	if len(c.Aliases) > 0 {
		aliasesString = fmt.Sprintf(" (aliases: %v)", c.Aliases)
	}
	usageString := ""
	if len(c.Usage) > 0 {
		usageString = fmt.Sprintf("\n\t└─Usage: %s", c.Usage)
	}
	text := au.Sprintf("%s%s: %s%s", au.Bold(c.Name), au.Gray(11, aliasesString), c.Help, usageString)
	return text
}

func getHelpString(c Cli, args []string) string {
	if len(args) >= 1 {		// If we have more than one args, find the command and print its help
		for _, command := range commandList {
			if command.Name == strings.ToLower(args[0]) {
				return command.description(c.au)
			}
		}
		return fmt.Sprintf("Unknown command: %s", args[0])
	}

	helpString := "Available commands:\n"
	for _, command := range commandList {
		helpString += command.description(c.au) + "\n"
	}
	return helpString
}

// This function is here so we can add the help command which refrences itself
func InitCommands() {
	commandList = append(commandList, command{
		Name:     "help",
		Aliases:  []string{"h", "?"},
		MinArgs:  0,
		Help:	  "Prints out help for all commands or a specified command",
		Usage:    "help [command]",
		Function: getHelpString,
	},)
}