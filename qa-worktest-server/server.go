package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	service := ":1919"
	if len(os.Args) > 1 {
		service = fmt.Sprintf(":%s", os.Args[1])
	}
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

	if _, err := conn.Write([]byte("READY>\r\n")); err != nil {
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

	log.Printf("got %d bytes from client: % 2x", n, buf[:n])

	if _, err = conn.Write([]byte(":)\r\n")); err != nil {
		log.Print(err)
	}

	deobfs := strings.Map(func(r rune) rune {
		return r + '+'
	}, string(buf))
	log.Printf("deo: %q", deobfs)
	data := make([]byte, 256)

	i := 0
	for i < len(deobfs) {
		n, err = base64.StdEncoding.Decode(data, []byte(deobfs[i:]))
		if err == nil {
			break
		}
		i++
	}
	log.Printf("Data: %q", data[:n])

	return
}
