# ssh public key encoding

go's crypto libraries are fairly comprehensive, but they don't cover the specific case of the "ssh one-line public key encoding".

This library provides methods to convert between `crypto/rsa.PublicKey` and the string ssh-rsa one-line encoding found in typical ssh `id_rsa.pub` files.  There is scaffolding in place to be able to extend it to other ciphers in the future, but that work is not yet done.

From the native structures such as `rsa.PublicKey`, it's fairly simple to use the existing crypto and encoding libraries to marshal into other formats such as PEM.

## Usage

```go

import "github.com/ianmcmahon/encoding_ssh"

pub_key, err := ssh.DecodePublicKey(string(bytes))

ssh_rsa_string, err := ssh.EncodePublicKey(user_key.PublicKey(), "user@host")

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


Here is another short program which reads $HOME/.ssh/id_rsa and outputs the public key in ssh-rsa one-line format.  This is equivalent to:

	ssh-keygen -y -f ~/.ssh/id_rsa

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
```