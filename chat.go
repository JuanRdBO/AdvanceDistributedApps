package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	conn net.Conn
	//nickname string
	ch chan<- string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func handleConnection(c net.Conn, msgchan chan<- string, addchan chan<- Client, rmchan chan<- net.Conn) {
	// ?
	ch := make(chan string)
	// ?
	msgs := make(chan string)
	// Create new channel for client
	addchan <- Client{c, ch}

	go func() {

		defer close(msgs)

		// Creates new reader
		bufc := bufio.NewReader(c)

		// Writes welcome message for first time
		c.Write([]byte("\033[1;30;41mWelcome to the fancy demo chat!. What is your nick? \033[0m\n"))
		
		// It reads the first message, which the client produces and assigns it to the nick.
		nick, _, err := bufc.ReadLine()

		// If error is produced while reading, return.
		if err != nil {
			return
		}

		// Assigns the previously assigned nick to "nickname"
		nickname := string(nick)

		// Writes to stdio from client the following message.
		c.Write([]byte("Welcome, " + nickname + "!\r\n\r\n"))

		// Assigns the following text to "msgs" to be displayed in stdio from server 
		// in the function "handleMessages".
		msgs <- "New user " + nickname + " has joined the chat room."

		// Reads all following messages from the user and displays it to stdio 
		// in the server side.
		for {
			// Reads from line.
			line, _, err := bufc.ReadLine()
			// If there's an error while reading, break.
			if err != nil {
				break
			}
			// Displays the message the user has written
			msgs <- nickname + ": " + string(line)
		}
		// When the loop ends, that means the user left the chat. Because the 
		// bufc.ReadLine above produced an error and broke the loop.
		msgs <- "User " + nickname + " left the chat room."
	}()

	// Loop to mantain connection open
	LOOP:
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					break LOOP
				}
				msgchan <- msg
			case msg := <-ch:
				_, err := c.Write([]byte(msg))
				if err != nil {
					break LOOP
				}
			}
		}

	c.Close()
	fmt.Printf("Connection from %v closed.\n", c.RemoteAddr())
	rmchan <- c
}

func handleMessages(c net.Conn, msgchan <-chan string, addchan <-chan Client, rmchan <-chan net.Conn) {
	clients := make(map[net.Conn]chan<- string)

	for {
		select {
		case msg := <-msgchan:
			
			// Trim suffix "\n" from message to be red as input.
			msg = strings.TrimSuffix(msg, "\n")
			// Prints to stdout. 
			fmt.Println("new message ", msg, " ( socket ", c.RemoteAddr(),")")

			// Renews the channel. Without it, the "Text to send" prompt, would not be displayed
			// after the second message.
			for _, ch := range clients {
					go func(mch chan<- string) { mch <- "\033[1;33;40m" + msg + "\033[m\r\n" }(ch)
			}

		case client := <-addchan:
			//fmt.Printf("New client: %v\n", client.conn)
			fmt.Println("New client ",c.RemoteAddr())
			clients[client.conn] = client.ch
		case conn := <-rmchan:
			//fmt.Printf("Client disconnects: %v\n", conn)
			fmt.Println("Cliet disconnects ",c.RemoteAddr())
			delete(clients, conn)
		}
	}
}

func main() {

	msgchan := make(chan string)
	addchan := make(chan Client)
	//rmchan := make(chan Client)
	rmchan := make(chan net.Conn)

	var string_server, port_server, ip_server, port_string_server string


	f, err := os.Open("configFile")
	if err != nil {
		fmt.Println("error opening file= ", err)
		os.Exit(1)
	}
	reader := bufio.NewReader(f)
	string_server, e := reader.ReadString('\n')

	counter := 1
	for e == nil {

		if counter == 1 {


			// Trim suffix \n from the string to add as port.
			string_server = strings.TrimSuffix(string_server, "\n")
			fmt.Println("\nstring: ", string_server, "and line: ", counter)
			ip_string_server := strings.Split(string_server, ":")
			ip_server, port_server = ip_string_server[0], ip_string_server[1]
			fmt.Println(port_server)
			fmt.Println("IP : ", ip_server, "\nPort : ", port_server, "\n")
			break
		}
		// if counter == 2 {

		// 	var string_message1, port_message1, ip_message1, port_string_message1 string

		// 	// Trim suffix \n from the string to add as port.
		// 	string_message1 = strings.TrimSuffix(string_message1, "\n")
		// 	fmt.Println("\nstring: ", string_message1, "and line: ", counter)
		// 	ip_string_message1 := strings.Split(string_server, ":")
		// 	ip_message1, port_message1 = ip_string_message1[0], ip_string_message1[1]
		// 	fmt.Println(port_message1)
		// 	fmt.Println("IP : ", ip_message1, "\nPort : ", port_message1, "\n")
		// 	break
		// }
		string_server, e = Readln(reader)
		counter = counter + 1

	}

	// connect to this socket
	//"127.0.0.1:6000"
	fmt.Println(string_server)

	//------------------------------Server------------------------------

	// To append ":" to "6000" to have valid port.
	var buf bytes.Buffer
	buf.WriteString(":")
	buf.WriteString(port_server)
	port_string_server = buf.String()

	// Trim suffix \n from the string to add as port.
	port_string_server = strings.TrimSuffix(port_string_server, "\n")

	// Initialize listener (SERVER)
	ln, err := net.Listen("tcp", port_string_server)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		
		go handleMessages(conn, msgchan, addchan, rmchan)

		go handleConnection(conn, msgchan, addchan, rmchan)
	}
}
