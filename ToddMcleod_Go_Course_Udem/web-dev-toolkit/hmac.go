package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
)

func main() {
	c := getCode("jswizzle@example.com")
	fmt.Println(c)
	d := getCode("jswizzle@example.comm")
	fmt.Println(d)

	isEqual := c == d
	fmt.Println("isEqual: ", isEqual) // output false, so we know user tampered with their stored hash.
}

func getCode(s string) string {
	h := hmac.New(sha256.New, []byte("ourkey"))
	_, err := io.WriteString(h, s)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
