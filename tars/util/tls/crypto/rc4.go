package crypto

import (
	"crypto/rc4"
)

func newRc4Cipher(key []byte) (cipher Cipher, err error) {
	k := len(key)
	if k < 1 || k > 256 {
		return nil, rc4.KeySizeError(k)
	}

	cipher = &rc4Cipher{
		key: key,
	}

	return
}

type rc4Cipher struct {
	key []byte
}

func (th *rc4Cipher) Encrypt(data []byte) ([]byte, error) {
	cipher, err := rc4.NewCipher(th.key)
	if err != nil {
		return nil, err
	}

	destData := make([]byte, len(data))
	cipher.XORKeyStream(destData, data)

	return destData, nil
}

func (th *rc4Cipher) Decrypt(data []byte) ([]byte, error) {
	cipher, err := rc4.NewCipher(th.key)
	if err != nil {
		return nil, err
	}

	destData := make([]byte, len(data))
	cipher.XORKeyStream(destData, data)

	return destData, nil
}
