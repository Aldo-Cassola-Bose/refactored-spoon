package main

import (
	"log"
	"net"
	"os"
	"time"
)

func main() {
	service := ":1919"
	listener, err := net.Listen("tcp", service)
	if err != nil {
		log.Printf("listen failed: %s", err)
		os.Exit(-1)
	}
	defer listener.Close()

	log.Printf("Listening on %s", listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("new connection err: %s", err)
		}

		log.Printf("new client at: %s", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		log.Printf("closing connection with %s", conn.RemoteAddr())
		conn.Close()
	}()

	buf := make([]byte, 256)

	_, err := conn.Write([]byte("READY>\r\n"))
	if err != nil {
		log.Print("error writing to client", conn.RemoteAddr())
		return
	}

	log.Print("Ready for input:")
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		log.Print("error reading from client ", err)
		return
	}

	log.Printf("got %d bytes from client: %q", n, buf[:n])
	_, err = conn.Write([]byte(":)\r\n"))
	if err != nil {
		log.Print(err)
	}
	return
}
