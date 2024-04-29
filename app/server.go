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

	var res []byte
	if header.path == "/" {
		res = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		res = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	_, err = conn.Write(res)
	if err != nil {
		fmt.Println(err)
	}
}
