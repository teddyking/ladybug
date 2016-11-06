package main

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/garden/client"
	"code.cloudfoundry.org/garden/client/connection"

	"github.com/jessevdk/go-flags"
	"github.com/teddyking/ladybug/commands"
	sys "github.com/teddyking/ladybug/system"
)

type command struct {
	name        string
	description string
	command     interface{}
}

type options struct {
	Depot string `short:"d" long:"depot" description:"Path to the garden depot dir" default:"/var/vcap/data/garden/depot"`
}

func main() {
	appOptions := &options{}
	parser := flags.NewParser(appOptions, flags.HelpFlag|flags.PassDoubleDash)

	parser.Parse()

	gardenClient := client.New(connection.New("tcp", "127.0.0.1:7777"))
	linuxHost := &sys.LinuxHost{DepotDir: appOptions.Depot, Proc: "/proc"}

	commands := []command{
		{"info", "Print info about garden and the host", &commands.Info{Client: gardenClient, Out: os.Stdout}},
		{"containers", "Print detailed info about containers on the host", &commands.Containers{Client: gardenClient, Out: os.Stdout, Host: linuxHost}},
	}

	for _, cmd := range commands {
		parser.AddCommand(cmd.name, cmd.description, "", cmd.command)
	}

	if _, err := parser.Parse(); err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err)
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
