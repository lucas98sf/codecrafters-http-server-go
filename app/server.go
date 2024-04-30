package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type HttpHeader struct {
	method  string
	path    string
	version string
}

type Header struct {
	key   string
	value string
}

const VERSION = "HTTP/1.1"
const CRLF = "\r\n"
const OK = VERSION + " 200 OK"
const NOT_FOUND = VERSION + " 404 Not Found"

const ContentPlain = "Content-Type: text/plain"

func contentLength(content string) string {
	return string([]byte("Content-Length: ")) + fmt.Sprint(len(content))
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	req := make([]byte, 1024)
	_, err := conn.Read(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	headers := strings.Split(string(req), "\r\n")
	http_header_str := strings.Split(headers[0], " ")
	http_header := HttpHeader{method: http_header_str[0], path: http_header_str[1], version: http_header_str[2]}

	path := strings.Split(string(http_header.path), "/")

	var res []byte
	if http_header.path == "/" {
		res = []byte(fmt.Sprintln(OK) + CRLF)
	} else if path[1] == "echo" {
		res = []byte(strings.Join([]string{OK, ContentPlain, contentLength(path[2]), CRLF + path[2]}, CRLF))
	} else if path[1] == "user-agent" {
		user_agent_str := strings.Split(headers[2], " ")
		user_agent_header := Header{key: user_agent_str[0], value: user_agent_str[1]}

		res = []byte(strings.Join([]string{OK,
			ContentPlain,
			contentLength(user_agent_header.value),
			CRLF + user_agent_header.value},
			CRLF))
	} else {
		res = []byte(fmt.Sprintln(NOT_FOUND) + CRLF)
	}

	_, err = conn.Write(res)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		return
	}

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}
