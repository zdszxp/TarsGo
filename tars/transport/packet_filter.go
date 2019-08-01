package transport

import (
	"github.com/TarsCloud/TarsGo/tars/util/tls/crypto"
)

type PacketFilter interface {
	Read([]byte) ([]byte, error)
	Write([]byte) ([]byte, error)
}

type EncryptionFilter interface {
	PacketFilter
	Encrypt(in []byte) (out []byte, err error)
	Decrypt(in []byte) (out []byte, err error)
}

type encryptionFilter struct {
	crypto.Cipher
}

func NewEncryptionFilter(typ string, privateKey []byte, publicKey []byte) (EncryptionFilter, error) {
	cipher, err := crypto.NewCipher(typ, privateKey, publicKey)
	if err != nil {
		return nil, err
	}
	return &encryptionFilter{
		Cipher: cipher,
	}, nil
}

func (th *encryptionFilter) Read(in []byte) ([]byte, error){
	return th.Decrypt(in)
}

func (th *encryptionFilter) Write(in []byte) ([]byte, error) {
	return th.Encrypt(in)
}