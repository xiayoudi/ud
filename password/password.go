// Copyright ©2026 xiayoudi. All rights reserved.
// Author: xiayoudi
// Email: ur@xiaud.com

package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/xiayoudi/ud"
	"golang.org/x/crypto/argon2"
)

const (
	defaultVariant = "argon2id"
	defaultVersion = argon2.Version
)

type Config struct {
	Time, Memory       uint32
	Threads            uint8
	KeyLength, SaltLen uint32
	MaxConcurrency     int
	WaitTimeout        time.Duration
}

type hasher struct {
	config    *Config
	semaphore chan struct{}
}

func DefaultConfig() *Config {
	return &Config{
		Time: 3, Memory: 1 << 15, Threads: 2,
		KeyLength: 32, SaltLen: 16,
		MaxConcurrency: 8, WaitTimeout: 5 * time.Second,
	}
}

func New(cfg *Config) *hasher {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return &hasher{
		config:    cfg,
		semaphore: make(chan struct{}, cfg.MaxConcurrency),
	}
}

var defaultHasher = New(nil)

func Hash(password string) (string, error) { return defaultHasher.Hash(password) }

func Validate(password, encodedHash string) (bool, error) {
	return defaultHasher.Validate(password, encodedHash)
}

func (h *hasher) Hash(password string) (string, error) {
	if password == "" {
		return "", ud.Err("password cannot be empty")
	}

	if err := h.acquire(); err != nil {
		return "", err
	}
	defer h.release()

	salt := make([]byte, h.config.SaltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", ud.Wrap(err, "failed to generate salt")
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.config.Time,
		h.config.Memory,
		h.config.Threads,
		h.config.KeyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%s$%s",
		defaultVariant, defaultVersion, h.config.Memory, h.config.Time, h.config.Threads, b64Salt, b64Hash), nil
}

func (h *hasher) Validate(password string, encodedHash string) (bool, error) {
	if password == "" || encodedHash == "" {
		return false, ud.Err("credentials missing")
	}

	p := strings.Split(encodedHash, "$")
	if len(p) != 6 || p[1] != defaultVariant {
		return false, ud.Err("invalid hash format or variant")
	}

	var v int
	var m, t uint32
	var th uint8

	if _, err := fmt.Sscanf(p[2], "v=%d", &v); err != nil || v != defaultVersion {
		return false, ud.Err("invalid version")
	}
	if _, err := fmt.Sscanf(p[3], "m=%d,t=%d,p=%d", &m, &t, &th); err != nil {
		return false, ud.Err("invalid parameters")
	}

	if err := h.acquire(); err != nil {
		return false, err
	}
	defer h.release()

	salt, err1 := base64.RawStdEncoding.DecodeString(p[4])
	actualHash, err2 := base64.RawStdEncoding.DecodeString(p[5])
	if err1 != nil || err2 != nil {
		return false, ud.Err("decode failed")
	}

	expectHash := argon2.IDKey([]byte(password), salt, t, m, th, uint32(len(actualHash)))
	return subtle.ConstantTimeCompare(expectHash, actualHash) == 1, nil
}

func (h *hasher) acquire() error {
	select {
	case h.semaphore <- struct{}{}:
		return nil
	case <-time.After(h.config.WaitTimeout):
		return ud.Err("server too busy")
	}
}

func (h *hasher) release() { <-h.semaphore }
