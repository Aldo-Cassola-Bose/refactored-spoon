package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	mathrand "math/rand"
	"net"
	"os"
	"strings"
	"time"
)

func usage(progname string) {
	fmt.Printf("%s {serverName portNumber}: run a qa test server\n", progname)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 {
		usage(os.Args[0])
	}
	service := os.Args[1] + ":" + os.Args[2]
	conn, err := net.Dial("tcp", service)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(-1)
	}

	log.Printf("connection opened with: %s", conn.RemoteAddr())
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("failed read from server: %s", err)
		os.Exit(-1)
	}

	if !bytes.Equal(buf[:n], []byte("READY>\r\n")) {
		log.Printf("unexpected data from server: %s", err)
		os.Exit(-1)
	}

	fmt.Printf(string(buf[:n]))
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	text := scan.Text()
	log.Print("sending data to server")

	encoded := strings.Split(
		base64.StdEncoding.EncodeToString([]byte(text)), "=")[0]
	block := make([]byte, 256)
	mathrand.Seed(time.Now().UnixNano())
	mathrand.Read(block)

	obfs := strings.Map(func(r rune) rune {
		return r - '+'
	}, encoded)

	copy(block[len(block)-len(obfs):], obfs)
	_, err = conn.Write([]byte(block))
	if err != nil {
		log.Printf("error writing to %s: %s", conn.RemoteAddr(), err)
		os.Exit(-1)
	}

	n, err = conn.Read(buf)
	if err != nil {
		log.Print(err)
		os.Exit(-1)
	}
	log.Printf("%s", buf[:n])
	log.Println("Done")
}
