package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	// Listen to incoming requests
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}

	// Close listener once main function terminates
	defer li.Close()

	// Infinite loop to continuously listen to incoming requests
	for {
		conn, err := li.Accept()
		if err != nil {
			log.Panic(err)
		}
		// Start an independent thread of control 'go' routine
		go handle(conn)
	}
}

// Once a request has been accepted, we handle that request.
func handle(conn net.Conn) {
	defer conn.Close()

	// read request
	request(conn)

	// write response
	respond(conn)
}

func request(conn net.Conn) {
	i := 0

	scanner := bufio.NewScanner(conn)
	// Loop over each section of the request
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)

		if i == 0 {
			fields := strings.Fields(ln)
			// request line - retrieve method portion of first line
			m := fields[0]
			// Print to standard out
			fmt.Println("***METHOD", m)
			// request line - retrieve URL portion of first line
			n := fields[1]
			// Print to standard out
			fmt.Println("***FIELDS", n)
		}

		if ln == "" {
			// headers are done
			break
		}

		i++
	}
}

func respond(conn net.Conn) {
	body := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Test</title></head><body><strong>Hello World</strong></body></html>`

	_, err := fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n")
	if err != nil {
		return
	}

	_, err2 := fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	if err2 != nil {
		return
	}
	_, err3 := fmt.Fprintf(conn, "Content-Type: text/html\r\n")
	if err3 != nil {
		return
	}
	_, err4 := fmt.Fprintf(conn, "\r\n")
	if err4 != nil {
		return
	}
	_, err5 := fmt.Fprintf(conn, body)
	if err5 != nil {
		return
	}
}
