package main

import (
	"fmt"
	"net"
)

type Command struct {
	Names		[]string
	Values		[]string
	function 	func(conn net.Conn, args []string)
}

type BotCmd struct {
	Commands	[]*Command
}

type Bot struct {
	Host     	string
	Port     	int
	CmdList 	BotCmd
}


func main() {
	commands := BotCmd {}
	cmd := Command {
		[]string{
			"test",
		},
		[]string{ },
		func(conn net.Conn, args []string) {
			fmt.Println("TEST!")
			_ = fmt.Sprintf("Args: %s", args)
		},
	}

	commands.Commands = append(commands.Commands, &cmd)

	bot := Bot {
		Host:    "127.0.0.1",
		Port:    6298,
		CmdList: commands,
	}

	bot.Start()
}
