package main

import (
	"github.com/ianmcmahon/encoding_ssh"

	"os"
	"fmt"
	"io/ioutil"
	"crypto/x509"
	"encoding/pem"
)

func main() {

	// read in public key from file
	bytes, err := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa.pub")
	if err != nil { fmt.Printf("%v\n", err); os.Exit(1) }

	// decode string ssh-rsa format to native type
	pub_key, err := ssh.DecodePublicKey(string(bytes))
	if err != nil { fmt.Printf("%v\n", err); os.Exit(1) }
	// pub_key is of type *rsa.PublicKey from 'crypto/rsa'

	// Marshal to ASN.1 DER encoding
	pkix, err := x509.MarshalPKIXPublicKey(pub_key)
	if err != nil { fmt.Printf("%v\n", err); os.Exit(1) }

	// Encode to PEM format
	pem := string(pem.EncodeToMemory(&pem.Block{
			Type: "RSA PUBLIC KEY",
			Bytes: pkix,
		}))

	fmt.Printf("%s", pem)
}
