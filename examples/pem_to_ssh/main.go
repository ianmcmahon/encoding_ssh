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

	// read in private key from file (private key is PEM encoded PKCS)
	bytes, err := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa")
	if err != nil { fmt.Printf("%v\n", err); os.Exit(1) }

	// decode PEM encoding to ANS.1 PKCS1 DER
	block, _ := pem.Decode(bytes)
	if block == nil { fmt.Printf("No Block found in keyfile\n"); os.Exit(1) }
	if block.Type != "RSA PRIVATE KEY" { fmt.Printf("Unsupported key type"); os.Exit(1) }

	// parse DER format to a native type
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	// encode the public key portion of the native key into ssh-rsa format
	// second parameter is the optional "comment" at the end of the string (usually 'user@host')
	ssh_rsa, err := ssh.EncodePublicKey(key.PublicKey, "")

	fmt.Printf("%s\n", ssh_rsa)
}
