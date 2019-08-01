package tls

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

var (
	ErrorInvalidKeyType = errors.New("invalid key type")
)

func readPrivateKey(privateKey string) (*ecdsa.PrivateKey, error) {
	var restPem []byte
	restPem = []byte(privateKey)
	for len(restPem) > 0 {
		var block *pem.Block
		block, restPem = pem.Decode(restPem)
		if block == nil {
			break
		}
		switch block.Type {
		case "EC PRIVATE KEY": // pkcs1
			return x509.ParseECPrivateKey(block.Bytes)
		case "PRIVATE KEY": // pkcs8
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err == nil {
				switch key.(type) {
				case *ecdsa.PrivateKey:
					return key.(*ecdsa.PrivateKey), err
				default:
					return nil, ErrorInvalidKeyType
				}
			} else {
				return nil, err
			}
		case "EC PARAMETERS":
			break
		default:
			return nil, ErrorInvalidKeyType
		}
	}
	return nil, errors.New("invalid pem")
}

func readPublicKey(publicKey string) (*ecdsa.PublicKey, error) {
	var restPem []byte
	restPem = []byte(publicKey)
	for len(restPem) > 0 {
		var block *pem.Block
		block, restPem = pem.Decode(restPem)
		if block == nil {
			break
		}
		switch block.Type {
		case "PUBLIC KEY": // pkcs1
			key, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err == nil {
				switch key.(type) {
				case *ecdsa.PublicKey:
					return key.(*ecdsa.PublicKey), nil
				default:
					return nil, ErrorInvalidKeyType
				}
			} else {
				return nil, err
			}
		default:
			return nil, ErrorInvalidKeyType
		}
	}
	return nil, errors.New("invalid pem")
}

func GenerateKey() (privateKeyBytes, publicKeyBytes []byte, err error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return
	}

	derStream, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return
	}

	keyOut := bytes.NewBuffer(nil)
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: derStream})
	privateKeyBytes = keyOut.Bytes()

	//public key
	publicKey := privateKey.PublicKey

	derStream, err = x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return
	}

	keyOut = bytes.NewBuffer(nil)
	pem.Encode(keyOut, &pem.Block{Type: "PUBLIC KEY", Bytes: derStream})
	publicKeyBytes = keyOut.Bytes()

	return
}

func GenerateKeyFile(privateKeyBytes, publicKeyBytes []byte) error {
	privateFile, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	defer privateFile.Close()

	_, err = privateFile.Write(privateKeyBytes)
	if err != nil {
		return err
	}

	publicFile, err := os.Create("public.pem")
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
