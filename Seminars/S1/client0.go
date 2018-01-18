package main

import (
	"bufio"
	"fmt"
	//"io"
	//"io/ioutil"
	"net"
	"os"
	"strings"
)

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

func main() {

	//dat, err := ioutil.ReadFile("configFile")
	//check(err)
	//fmt.Print(string(dat))
	string := ""

	f, err := os.Open("configFile")
	if err != nil {
		fmt.Println("error opening file= ", err)
		os.Exit(1)
	}
	reader := bufio.NewReader(f)
	string, e := reader.ReadString('\n')

	counter := 1
	for e == nil {

		if counter == 3 {
			fmt.Println("\nstring: ", string, "and line: ", counter)
			ip_string := strings.Split(string, ":")
			ip, port := ip_string[0], ip_string[1]
			fmt.Println("IP : ", ip, "\nPort : ", port, "\n")
			break
		}
		string, e = Readln(reader)
		counter = counter + 1

	}

	// connect to this socket
	//"127.0.0.1:6000"
	fmt.Println(string)

	//----------------------------END--------------------------------

	conn, _ := net.Dial("tcp", string)
	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// send to socket
		fmt.Fprintf(conn, text)
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
