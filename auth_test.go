package main

import "testing"

func TestAuthKey(t *testing.T) {
	auth := Auth{Secret: "Ceci n'est pas un string"}
	hash := "5cd95b79ec85715d3cff2679c36d88282a80cec0"
	key := auth.Key()
	if key != hash {
		t.Fatalf("Expected: %s, Got: %s", hash, key)
	}
}
