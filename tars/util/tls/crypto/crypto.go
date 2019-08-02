package crypto

import (
	"errors"
	"fmt"
)

const (
	CRYPTO_TYPE_AES = "AES"
	CRYPTO_TYPE_RC4 = "RC4"
	CRYPTO_TYPE_RSA = "RSA"
)

type Cipher interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

func NewCipher(typ string, key []byte, keys ...[]byte) (Cipher, error) {
	switch typ {
	case CRYPTO_TYPE_AES:
		return newAesCipher(key)
	case CRYPTO_TYPE_RC4:
		return newRc4Cipher(key)
	case CRYPTO_TYPE_RSA:
		if len(keys) < 1 {
			return nil, errors.New("NewCipher RSA Error: privateKey and publicKey must be set")
		}
		return newRsaCipher(key, keys[0])
	default:
		return nil, errors.New(fmt.Sprintf("crypto::NewCipher: not support type[%v]", typ))
	}
}
