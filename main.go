package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/efy/gorecall/cmd"
)

var helpTemplate = `
gorecall is a server application and set of tools for managing web links.

Usage:

	gorecall command [arguments]

The commands are:
{{ range . }}
	{{ .Name }}		{{ .Short }}{{ end }}

Use "gorecall [command] help" for more information about a command

`

func main() {
	// print help template when no commands are specified
	// or the the command is "help"
	if len(os.Args) < 2 || os.Args[1] == "help" {
		tmpl := template.Must(template.New("help").Parse(helpTemplate))
		tmpl.ExecuteTemplate(os.Stderr, "help", cmd.Commands)
		os.Exit(2)
	}

	// Run command
	for _, c := range cmd.Commands {
		if c.Name() == os.Args[1] && c.Runnable() {
			c.Run(&c, os.Args[2:])
			return
		}
	}

	// Handle invalid / unknown command
	fmt.Fprintf(os.Stderr, "gorecall: unknown subcommand %q\nRun 'gorecall help' for usage.\n", os.Args[1])
}
