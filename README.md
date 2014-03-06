# ssh public key encoding

go's crypto libraries are fairly comprehensive, but they don't cover the specific case of the "ssh one-line public key encoding".

This library provides methods to convert between `crypto/rsa.PublicKey` and the string ssh-rsa one-line encoding found in typical ssh `id_rsa.pub` files.  There is scaffolding in place to be able to extend it to other ciphers in the future, but that work is not yet done.

From the native structures such as `rsa.PublicKey`, it's fairly simple to use the existing crypto and encoding libraries to marshal into other formats such as PEM.

## Usage

```go
	import "github.com/ianmcmahon/encoding_ssh"

	pub_key, err := ssh.DecodePublicKey(string(bytes))
```

Here's a short program which reads $HOME/.ssh/id_rsa.pub and outputs the public key in PKCS8 format.  This is equivalent to:

	ssh-keygen -f $HOME/.ssh/id_rsa.pub -e -m pkcs8

```go
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
	// pub_key is of type *rsa.PublicKey

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
```
