package main

import (
	"fmt"
	"net"
)

type Command struct {
	Names		[]string
	Values		[]string
	function 	func(conn net.Conn)
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
	cmds := BotCmd {}

	cmd := Command {
		[]string{
			"test",
		},
		[]string{ },
		func(conn net.Conn) {
			fmt.Println("TEST!")
		},
	}

	cmds.Commands = append(cmds.Commands, &cmd)

	bot := Bot {
		Host: "127.0.0.1",
		Port: 6298,
		CmdList: cmds,
	}

	bot.Start()
}
