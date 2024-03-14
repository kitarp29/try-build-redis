package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")
	common := make(map[string]string)
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
		go handleConnection(conn, common)
	}

}

func handleConnection(conn net.Conn, common map[string]string) {
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
		} else if strings.Contains(str, "set") {
			split := strings.Split(str, "\r\n")
			fmt.Println(split[4], split[6])
			common[split[len(split)-4]] = split[len(split)-2]
			response := "+OK\r\n"
			_, err = conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error writing data to connection: ", err.Error())
				os.Exit(1)
			}
		} else if strings.Contains(str, "get") {
			split := strings.Split(str, "\r\n")
			response := "+" + common[split[len(split)-2]] + "\r\n"
			_, err = conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error writing data to connection: ", err.Error())
				os.Exit(1)
			}
		} else {
			split := strings.Split(str, "\r\n")
			response := "+" + split[len(split)-2] + "\r\n"
			_, err = conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error writing data to connection: ", err.Error())
				os.Exit(1)
			}
		}
	}
}
