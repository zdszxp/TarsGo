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

func NewCipher(typ string, privateKey []byte, publicKey []byte) (Cipher, error) {
	switch typ {
	case CRYPTO_TYPE_AES:
		return newAesCipher(privateKey)
	case CRYPTO_TYPE_RC4:
		return newRc4Cipher(privateKey)
	case CRYPTO_TYPE_RSA:
		return newRsaCipher(privateKey, publicKey)
	default:
		return nil, errors.New(fmt.Sprintf("crypto::NewCipher: not support type[%v]", typ))
	}
}
