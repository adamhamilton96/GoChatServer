package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

type user struct {
	name string
	conn net.Conn
}

var clients []user

func handleInput(client *user) {
	ln := ""
	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		ln = scanner.Text()
		if ln == "exit" {
			os.Exit(1)
		} else if client.name == "" {
			client.name = ln
		}
		sendAllOther(client.conn, ln)
	}
}

func sendAllOther(conn net.Conn, ln string) {
	name := ""
	for i := 0; i < len(clients); i++ {
		if clients[i].conn == conn {
			name = clients[i].name
		}
	}
	for i := 0; i < len(clients); i++ {
		if clients[i].conn != conn {
			io.WriteString(clients[i].conn, name+": "+ln+"\n")
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		name := ""
		client := user{name, conn}
		clients = append(clients, client)

		fmt.Println("New user joined, number of users = " + strconv.Itoa(len(clients)))
		for i := 0; i < len(clients); i++ {
			go handleInput(&clients[i])
		}
	}
}
