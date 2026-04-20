package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()
// 	fmt.Printf("The connection accepted. %s \n", conn.RemoteAddr())

// 	buf := make([]byte, 1024)
// 	dialConn, errDial := net.Dial("tcp", ":2001")

// 	for {
// 		n, err := conn.Read(buf)

// 		if err != nil{
// 			fmt.Printf("connection %s Closed", conn.RemoteAddr())
// 			dialConn.Close()
// 			fmt.Printf("Server Connection %s Closed", dialConn.RemoteAddr())
// 			return
// 		}
// 		if errDial != nil {
// 			fmt.Println("Error on dialing to server ...")
// 			conn.Close()
// 			return
// 		}

// 		dialConn.Write(buf[:n])

// 		fmt.Println(string(buf[:n]))
// 	}

// }



func proxyConnection(clientConn net.Conn) {
	defer clientConn.Close()
	serverConn, err := net.Dial("tcp", ":2001")

	if err!=nil{
		fmt.Println("Server connection problem", err)
		return
	}
	fmt.Println("Server Connection accepted", serverConn.RemoteAddr())
	defer serverConn.Close()

	go func(){
		tee := io.TeeReader(serverConn, os.Stdout)
		io.Copy(clientConn, tee)

	}()
	
	tee := io.TeeReader(clientConn, os.Stdout)
	io.Copy(serverConn, tee)






}

func main() {

	address := ":2000"

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Listen error", err)
	}

	defer listener.Close()

	log.Printf("Server is ready on the Port %s", address)
	// serverChan := make(chan net.Conn)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("connection error")
		}
		fmt.Println("Client Connection accepted", conn.RemoteAddr())

		// Each connection gets its own goroutine so one blocked client does not
		// stop the server from accepting and proxying other connections.
		go proxyConnection(conn)
		// go handleConnection(conn)

		// go clientConnection(conn,  serverChan)
		// go serverConnection(conn,  serverChan)

	}



}
