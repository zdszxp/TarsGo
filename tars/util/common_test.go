package util

import "testing"

func TestEqualbytes(t *testing.T) {
	bs0 := RandomCreateBytes(16)
	bs1 := RandomCreateBytes(16)

	t.Log(string(bs0), string(bs1))
	if Equalbytes(bs0, bs1) {
		t.Fatal()
	}

	if Equalbytes(bs0, bs0) == false {
		t.Fatal()
	}
}
