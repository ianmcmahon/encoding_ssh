package ssh

import (
	"fmt"
	"bytes"
	"strings"
	"encoding/binary"
	"encoding/base64"
	"crypto/rsa"
	"math/big"
)

// ssh one-line format (for lack of a better term) consists of three text fields: { key_type, data, comment }
// data is base64 encoded binary which consists of tuples of length (4 bytes) and data of the length described previously.
// For RSA keys, there should be three tuples which should be:  { key_type, public_exponent, modulus } 

func EncodePublicKey(key interface{}, comment string) (string, error) {
	if rsaKey, ok := key.(rsa.PublicKey); ok {
		key_type := "ssh-rsa"

		modulus_bytes := rsaKey.N.Bytes()

		buf := new(bytes.Buffer)

		var data = []interface{} {
			uint32(len(key_type)),
			[]byte(key_type),
			uint32(binary.Size(uint32(rsaKey.E))),
			uint32(rsaKey.E),
			uint32(binary.Size(modulus_bytes)),
			modulus_bytes,
		}

		for _, v := range data {
			err := binary.Write(buf, binary.BigEndian, v)
			if err != nil { return "", err }
		}

		return fmt.Sprintf("%s %s %s", key_type, base64.StdEncoding.EncodeToString(buf.Bytes()), comment), nil
	}

	return "", fmt.Errorf("Unknown key type: %T\n", key)
}

func readLength(data []byte) ([]byte, uint32, error) {
	l_buf := data[0:4]

	buf := bytes.NewBuffer(l_buf)

	var length uint32

	err := binary.Read(buf, binary.BigEndian, &length)
	if err != nil { return nil, 0, err }

	return data[4:], length, nil
}

func readBigInt(data []byte, length uint32) ([]byte, *big.Int, error) {
	var bigint = new(big.Int)
	bigint.SetBytes(data[0:length])
	return data[length:], bigint, nil
}

func getRsaValues(data []byte) (format string, e *big.Int, n *big.Int, err error) {
	data, length, err := readLength(data)
	if err != nil { return }

	format = string(data[0:length]); data = data[length:]

	data, length, err = readLength(data)
	if err != nil { return }

	data, e, err = readBigInt(data, length)
	if err != nil { return }

	data, length, err = readLength(data)
	if err != nil { return }

	data, n, err = readBigInt(data, length)
	if err != nil { return }

	return
}

func DecodePublicKey(str string) (interface{}, error) {
	// comes in as a three part string
	// split into component parts

	tokens := strings.Split(str, " ")

	if len(tokens) < 2 { return nil, fmt.Errorf("Invalid key format; must contain at least two fields (keytype data [comment])") }

	key_type := tokens[0]
	data, err := base64.StdEncoding.DecodeString(tokens[1])
	if err != nil { return nil, err }

	format, e, n, err := getRsaValues(data)

	if format != key_type { return nil, fmt.Errorf("Key type said %s, but encoded format said %s.  These should match!", key_type, format) }

	pubKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return pubKey, nil
}