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

const VERSION = "HTTP/1.1"
const CRLF = "\r\n"
const OK = VERSION + " 200 OK"
const NOT_FOUND = VERSION + " 404 Not Found"

const ContentPlain = "Content-Type: text/plain"

func contentLength (content string) string {
	return string([]byte("Content-Length: ")) + fmt.Sprint(len(content))
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	req := make([]byte, 1024)
	_, err = conn.Read(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	str_req := strings.Split(string(req), " ")
	header := HttpHeader{method: str_req[0], path: str_req[1], version: str_req[2]}
	path := strings.Split(string(header.path), "/")

	var res []byte
	if header.path == "/" {
		res = []byte(fmt.Sprintln(OK) + CRLF)
	} else if path[1] == "echo" {
		res = []byte(strings.Join([]string{OK, ContentPlain, contentLength(path[2]), CRLF + path[2]}, CRLF))
	} else {
		res = []byte(fmt.Sprintln(NOT_FOUND) + CRLF)
	}

	fmt.Println(string(res))

	_, err = conn.Write(res)
	if err != nil {
		fmt.Println(err)
	}
}
