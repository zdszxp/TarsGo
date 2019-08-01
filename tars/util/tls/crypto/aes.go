package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func newAesCipher(key []byte) (cipher Cipher, err error) {
	keyLen := len(key)
	if keyLen < 16 {
		err = fmt.Errorf("The length of res key shall not be less than 16")
		return
	}

	cipher = &aesCipher{
		key: key,
	}

	return
}

type aesCipher struct {
	key []byte
}

func (e *aesCipher) getKey() []byte {
	keyLen := len(e.key)
	if keyLen < 16 {
		panic("The length of res key shall not be less than 16")
	}

	if keyLen >= 32 {
		return e.key[:32]
	}
	if keyLen >= 24 {
		return e.key[:24]
	}

	return e.key[:16]
}

func (e *aesCipher) Encrypt(data []byte) ([]byte, error) {
	key := e.getKey()
	var iv = []byte(key)[:aes.BlockSize]

	encrypted := make([]byte, len(data))
	encrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	
	aesEncrypter := cipher.NewCFBEncrypter(encrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, data)
	return encrypted, nil
}

func (e *aesCipher) Decrypt(data []byte) ([]byte, error) {
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	key := e.getKey()
	var iv = []byte(key)[:aes.BlockSize]

	decrypted := make([]byte, len(data))
	decrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	
	aesDecrypter := cipher.NewCFBDecrypter(decrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, data)
	return decrypted, nil
}
