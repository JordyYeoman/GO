package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello sir")
	resp, err := http.Get("http://google.com")

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	lw := LogWriter{}

	io.Copy(lw, resp.Body)
	//io.Copy(os.Stdout, resp.Body)
}

// Custom writer
type LogWriter struct{}

func (LogWriter) Write(bs []byte) (int, error) {
	fmt.Println(string(bs))
	fmt.Println("Just wrote out N number of bytes: ", len(bs))
	return len(bs), nil
}
