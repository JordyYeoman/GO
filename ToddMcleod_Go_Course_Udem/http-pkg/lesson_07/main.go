package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	fmt.Println("Full send")

	// Basic tcp server
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("Yayaya - cannot accept incoming request")
		}
		// Run concurrently
		go serve(c)
	}

	return
}

func serve(c net.Conn) {
	defer c.Close()
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if ln == "" {
			// when ln is empty, header is done
			fmt.Println("THIS IS END OF THE NET CONN")
			break
		}
	}

	//
	body := "THIS IS THE BODDYYYY"
	io.WriteString(c, "HTTP/1.1 200 OK\r\n")
	io.WriteString(c, "\r\n")
	io.WriteString(c, body)
}
