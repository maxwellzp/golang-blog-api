package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	b := make([]byte, 32) // 256 bits
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	// this will produce a string of ~44 characters
	fmt.Println("JWT secret key:", base64.StdEncoding.EncodeToString(b))
}
