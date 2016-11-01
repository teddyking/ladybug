package main

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/garden/client"
	"code.cloudfoundry.org/garden/client/connection"

	"github.com/jessevdk/go-flags"
	"github.com/teddyking/ladybug/commands"
)

type command struct {
	name        string
	description string
	command     interface{}
}

func main() {
	gardenClient := client.New(connection.New("tcp", "127.0.0.1:7777"))
	parser := flags.NewParser(nil, flags.HelpFlag|flags.PassDoubleDash)

	commands := []command{
		{"info", "Prints info about garden and the host", &commands.Info{Client: gardenClient, Out: os.Stdout, Err: os.Stderr}},
	}

	for _, cmd := range commands {
		parser.AddCommand(cmd.name, cmd.description, "", cmd.command)
	}

	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err)
	}
}
