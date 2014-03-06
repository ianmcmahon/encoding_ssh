# ssh public key encoding

go's crypto libraries are fairly comprehensive, but they don't cover the specific case of the "ssh one-line public key encoding".

This library provides methods to convert between `crypto/rsa.PublicKey` and the string ssh-rsa one-line encoding found in typical ssh `id_rsa.pub` files.  There is scaffolding in place to be able to extend it to other ciphers in the future, but that work is not yet done.

From the native structures such as `rsa.PublicKey`, it's fairly simple to use the existing crypto and encoding libraries to marshal into other formats such as PEM.

## Usage

```go
package main

import (
	
)
```
