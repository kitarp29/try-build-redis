package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("Logs from your program will appear here!")
	// common is a map that stores key-value pairs, where the keys are of type string
	// and the values are of type commonValue.
	var common = make(map[string]string)
	// commonValue represents a value in the common map, which consists of a string and an integer.

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
				common[commands[i+1]] = commands[i+2]

				_, err = conn.Write([]byte("+OK\r\n"))
				if err != nil {
					fmt.Println("Error writing data to connection: ", err.Error())
				}
				i += 2
				break // Skip to the next iteration
			case "get":
				if val, ok := common[commands[i+1]]; ok {
					response := "+" + val + "\r\n"
					i++
					_, err = conn.Write([]byte(response))
					if err != nil {
						fmt.Println("Error writing data to connection: ", err.Error())
					}
					i++
					break
				} else {
					_, err = conn.Write([]byte("$-1\r\n"))
				} // Skip to the next iteration
			case "px":
				t, err := strconv.Atoi(commands[i+1])
				if err != nil {
					fmt.Println("Error converting string to integer: ", err.Error())
					continue
				}
				go func(key string, t int) {
					time.Sleep(time.Duration(t) * time.Millisecond)
					delete(common, key)
				}(commands[i-2], t)
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
