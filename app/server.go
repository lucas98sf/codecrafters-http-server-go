package main

import (
	"flag"
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

const (
	Version      = "HTTP/1.1"
	NextLine     = "\r\n"
	Ok           = Version + " 200 OK"
	NotFound     = Version + " 404 Not Found"
	ContentPlain = "Content-Type: text/plain"
)

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
		res = []byte(fmt.Sprintln(Ok) + NextLine)
	} else if path[1] == "echo" {
		res = []byte(strings.Join([]string{Ok, ContentPlain, contentLength(path[2]), NextLine + path[2]}, NextLine))
	} else if path[1] == "user-agent" {
		user_agent_str := strings.Split(headers[2], " ")
		user_agent_header := Header{key: user_agent_str[0], value: user_agent_str[1]}

		res = []byte(strings.Join([]string{Ok,
			ContentPlain,
			contentLength(user_agent_header.value),
			NextLine + user_agent_header.value},
			NextLine))
	} else if path[1] == "files" {
		if http_header.method == "GET" {
			content, err := os.ReadFile(directory + "/" + path[2])
			if err != nil {
				res = []byte(fmt.Sprintln(NotFound) + NextLine)
				_, err = conn.Write(res)
				if err != nil {
					fmt.Println(err)
					return
				}
				return
			}
			res = []byte(strings.Join([]string{Ok,
				"Content-Type: application/octet-stream",
				contentLength(string(content)),
				NextLine + string(content)},
				NextLine))
		}
	} else {
		res = []byte(fmt.Sprintln(NotFound) + NextLine)
	}

	_, err = conn.Write(res)
	if err != nil {
		fmt.Println(err)
		return
	}
}

var directory string

func main() {
	flag.StringVar(&directory, "directory", "", "")
	flag.Parse()
	fmt.Println("Directory:", directory)
	if directory != "" {
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			fmt.Println("Directory does not exist", err)
			os.Exit(1)
		}
	}

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
