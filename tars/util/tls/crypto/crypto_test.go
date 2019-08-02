package crypto

import (
	"testing"
	"bytes"
	"github.com/TarsCloud/TarsGo/tars/util"
)

func TestAESCipher(t *testing.T){
	typ := CRYPTO_TYPE_AES
	key := util.RandomCreateBytes(16)

	testCipher(t, typ, key)
}

func testCipher(t *testing.T, typ string, keys... []byte) {
	cipher, err := NewCipher(typ, keys[0], nil)
	if err != nil {
		t.Fatal(err)
	}

	//test multi times
	for i := 10; i != 0 ; i-- {
		data := util.RandomCreateBytes(512)

		encryptedData, err := cipher.Encrypt(data)
		if err != nil {
			t.Fatal(err)
		}
	
		decryptedData, err := cipher.Decrypt(encryptedData)
		if err != nil {
			t.Fatal(err)
		}
	
		if !bytes.Equal(data, decryptedData) {
			t.Fatal()
		}
	}
}