// Copyright ©2026 xiayoudi. All rights reserved.
// Author: xiayoudi
// Email: ur@xiaud.com

package aes

import "testing"

func TestAES(t *testing.T) {
	// openssl rand -base64 32 | cut -c1-32
	key := []byte("1vNdT06KieH20VGpFNcbSp0DowSiZby7")
	Init(key)

	plaintext := []byte("hello ud 2026")
	aad := []byte("ur")

	ct, err := Encrypt(plaintext, aad)
	if err != nil {
		t.Fatal(err)
	}

	pt, err := Decrypt(ct, aad)
	if err != nil {
		t.Fatal(err)
	}

	if string(pt) != string(plaintext) {
		t.Error("plaintext mismatch")
	}
}
