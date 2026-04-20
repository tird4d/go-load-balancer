package main

import (
	"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("The connection accepted. %s", conn.RemoteAddr())

	buf := make([]byte, 1024)
	for {
		conn.Read(buf)
		fmt.Println(string(buf))
	}

}

func main() {
	//ghazalmajesty

	address := ":8080"

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Listen error", err)
	}

	defer listener.Close()

	log.Printf("Server is ready on the Port %s", address)

	conn, err := listener.Accept()

	if err != nil {
		log.Fatal("connection error")
	}

	handleConnection(conn)

}
