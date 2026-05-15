// Copyright ©2026 xiayoudi. All rights reserved.
// Author: xiayoudi
// Email: ur@xiaud.com

package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"sync"

	"github.com/xiayoudi/ud"
)

var (
	global *aesgcm
	once   sync.Once
)

type aesgcm struct {
	gcm cipher.AEAD
}

func New(key []byte) (*aesgcm, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &aesgcm{gcm: gcm}, nil
}

func Init(key []byte) {
	once.Do(func() {
		instance, err := New(key)
		if err != nil {
			panic("uaes init failed: " + err.Error())
		}
		global = instance
	})
}

func Global() *aesgcm {
	if global == nil {
		panic("uaes: global instance not initialized, call Init() first")
	}
	return global
}

func (a *aesgcm) encrypt(plaintext, aad []byte) ([]byte, error) {
	nonce := make([]byte, a.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return a.gcm.Seal(nonce, nonce, plaintext, aad), nil
}

func (a *aesgcm) decrypt(data, aad []byte) ([]byte, error) {
	ns := a.gcm.NonceSize()
	if len(data) < ns {
		return nil, ud.Err("ciphertext too short")
	}

	nonce := data[:ns]
	ciphertext := data[ns:]

	return a.gcm.Open(nil, nonce, ciphertext, aad)
}

func Encrypt(plaintext, aad []byte) ([]byte, error) { return Global().encrypt(plaintext, aad) }

func Decrypt(data, aad []byte) ([]byte, error) { return Global().decrypt(data, aad) }
