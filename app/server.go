package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	data := make([]byte, 1024)
	for {
		s, err := conn.Read(data)
		if err != nil {
			fmt.Println("Error reading data from connection: ", err.Error())
			os.Exit(1)
		}
		str := (string(data[:s]))
		if strings.Contains(str, "ping") {
			response := "+PONG\r\n"
			_, err = conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error writing data to connection: ", err.Error())
				os.Exit(1)
			}
		} else {
			split := strings.Split(str, "\r\n")
			// fmt.Println(split, len(split))
			// fmt.Printl=n(split[len(split)-2])
			response := "+" + split[len(split)-2] + "\r\n"
			// fmt.Println(response)
			_, err = conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error writing data to connection: ", err.Error())
				os.Exit(1)
			}
		}
	}
}
