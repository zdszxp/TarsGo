package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func newRsaCipher(privateKeyBytes, publicKeyBytes []byte) (cipher Cipher, err error) {
	if privateKeyBytes == nil && publicKeyBytes == nil {
		return nil, errors.New("crypto::newRsaCipher:: invalid rsa key")
	}

	rsaCipher := &rsaCipher{}

	cipher = rsaCipher

	//maybe just encrypt
	if publicKeyBytes != nil {
		block, _ := pem.Decode(publicKeyBytes)
		if block == nil {
			err = errors.New("crypto::newRsaCipher: public key error")
			return
		}

		var pub interface{}
		pub, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return
		}

		var ok bool
		rsaCipher.publicKey, ok = pub.(*rsa.PublicKey)
		if !ok {
			err = errors.New("crypto::newRsaCipher: public key error")
			return
		}
	}

	//maybe just decrypt
	if privateKeyBytes != nil {
		block, _ := pem.Decode(privateKeyBytes)
		if block == nil {
			err = errors.New("crypto::newRsaCipher: private key error")
			return
		}

		rsaCipher.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return
		}
	}

	return
}

type rsaCipher struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (th *rsaCipher) Encrypt(data []byte) ([]byte, error) {
	if th.publicKey == nil {
		return nil, errors.New("crypto::Encrypt:: invalid rsa publicKey")
	}

	return rsa.EncryptPKCS1v15(rand.Reader, th.publicKey, data)
}

func (th *rsaCipher) Decrypt(data []byte) ([]byte, error) {
	if th.privateKey == nil {
		return nil, errors.New("crypto::Encrypt:: invalid rsa privateKey")
	}

	return rsa.DecryptPKCS1v15(rand.Reader, th.privateKey, data)
}

func GenerateRSAKey() (privateKeyBytes, publicKeyBytes []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return
	}

	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	if derStream == nil {
		return
	}

	keyOut := bytes.NewBuffer(nil)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: derStream})
	privateKeyBytes = keyOut.Bytes()

	//public key
	derStream, err = x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return
	}

	keyOut = bytes.NewBuffer(nil)
	pem.Encode(keyOut, &pem.Block{Type: "PUBLIC KEY", Bytes: derStream})
	publicKeyBytes = keyOut.Bytes()

	return
}

func GenerateRsaKeyFile(privateKeyBytes, publicKeyBytes []byte) error {
	privateFile, err := os.Create("private_rsa.pem")
	if err != nil {
		return err
	}
	defer privateFile.Close()

	_, err = privateFile.Write(privateKeyBytes)
	if err != nil {
		return err
	}

	publicFile, err := os.Create("public_rsa.pem")
	if err != nil {
		return err
	}
	defer publicFile.Close()

	_, err = publicFile.Write(publicKeyBytes)
	if err != nil {
		return err
	}

	return nil
}
