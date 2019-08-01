package util

import "testing"

func TestRandomCreateBytes(t *testing.T) {
	bs0 := RandomCreateBytes(16)
	bs1 := RandomCreateBytes(16)

	t.Log(string(bs0), string(bs1))
	if string(bs0) == string(bs1) {
		t.Fatal("random same bytes")
	}

	bs0 = RandomCreateBytes(4, []byte(`a`)...)

	if string(bs0) != "aaaa" {
		t.Fatal("wrong result")
	}
}
