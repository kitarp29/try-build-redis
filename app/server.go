package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
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
func parseRedisCommand(input []byte) []string {
	rawCommand := string(input)
	commands := strings.Split(rawCommand, "\r\n")
	finalCommands := []string{}
	// parse arrays
	if strings.HasPrefix(commands[0], "*") {
		_, err := strconv.Atoi(commands[0][1:])
		if err != nil {
			return []string{"Encountered error while parsing *"}
		}
		checkLengthFlag := false
		for _, v := range commands[1:] {
			if strings.HasPrefix(v, "$") {
				_, err := strconv.Atoi(v[1:])
				if err != nil {
					return []string{"Encountered error while parsing $"}
				}
				checkLengthFlag = true
			} else if checkLengthFlag {
				checkLengthFlag = false
				finalCommands = append(finalCommands, v)
			}
		}
	}
	return finalCommands
}

func handleConnection(conn net.Conn, common map[string]string) {
	defer conn.Close()

	data := make([]byte, 1024)
	for {
		_, err := conn.Read(data)
		commands := parseRedisCommand(data)
		if err != nil {
			fmt.Println("Error reading data from connection: ", err.Error())
			break
		}

		fmt.Println("Received command: ", commands)
		for i, command := range commands {
			switch command {
			case "ping":
				_, err = conn.Write([]byte("+PONG\r\n"))
				if err != nil {
					fmt.Println("Error writing data to connection: ", err.Error())
				}
			case "set":
				if i+1 < len(commands) {
					common[commands[i+1]] = commands[i+2]
					_, err = conn.Write([]byte("+OK\r\n"))
					if err != nil {
						fmt.Println("Error writing data to connection: ", err.Error())
					}
					i += 2
					break // Skip to the next iteration
				}
			case "get":
				if i+1 < len(commands) {
					response := "+" + common[commands[i+1]] + "\r\n"
					i++
					_, err = conn.Write([]byte(response))
					if err != nil {
						fmt.Println("Error writing data to connection: ", err.Error())
					}
					i++
					break // Skip to the next iteration
				}
			case "echo":
				if i+1 < len(commands) {
					response := "+" + commands[i+1] + "\r\n"
					_, err = conn.Write([]byte(response))
					if err != nil {
						fmt.Println("Error writing data to connection: ", err.Error())
					}
					i++
					break // Skip to the next iteration
				}
			default:
				continue
			}
		}
	}
}
