package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	CRLF := "\r\n"

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		req := make([]byte, 1024)
		conn.Read(req)
		url := strings.Split(string(req), CRLF)[0]

		var response string
		if !strings.HasPrefix(url, "GET / HTTP/1.1") {
			response = "HTTP/1.1 404 Not Found" + CRLF + CRLF
		} else {
			response = "HTTP/1.1 200 OK" + CRLF + CRLF
		}

		conn.Write([]byte(response))
		conn.Close()
	}
}
