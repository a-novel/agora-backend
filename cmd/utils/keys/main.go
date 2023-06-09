package main

// Generates a private key in PEM format, for testing purposes.

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
)

func main() {
	_, pk, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err.Error())
	}

	mpk, err := x509.MarshalPKCS8PrivateKey(pk)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Raw private key:")

	var displayBytes []string
	for _, b := range mpk {
		displayBytes = append(displayBytes, fmt.Sprintf("%v", b))
	}
	fmt.Println("[", strings.Join(displayBytes, ", "), "]")

	fmt.Println("\nFile version:")
	err = pem.Encode(os.Stdout, &pem.Block{Type: "PRIVATE KEY", Bytes: mpk})
	if err != nil {
		panic(err.Error())
	}
}
