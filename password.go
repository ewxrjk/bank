package bank

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"strconv"
	"strings"
)

// ErrPasswordMismatch is returned by VerifyPassword for incorrect passwords.
var ErrPasswordMismatch = errors.New("password mismatch")

// VerifyPassword verifies a password against a stored version.
func VerifyPassword(stored, password string) (err error) {
	bits := strings.Split(stored, ",")
	alg := bits[0]
	if alg == "scrypt" {
		var n, r, p int
		if n, err = strconv.Atoi(bits[1]); err != nil {
			return
		}
		if r, err = strconv.Atoi(bits[2]); err != nil {
			return
		}
		if p, err = strconv.Atoi(bits[3]); err != nil {
			return
		}
		var salt, hash []byte
		if salt, err = base64.RawStdEncoding.DecodeString(bits[4]); err != nil {
			return
		}
		if hash, err = base64.RawStdEncoding.DecodeString(bits[5]); err != nil {
			return
		}
		var xhash []byte
		if xhash, err = scrypt.Key([]byte(password), salt, n, r, p, 16); err != nil {
			return
		}
		if string(xhash) != string(hash) {
			err = ErrPasswordMismatch
			return
		}
	} else {
		return errors.New("unrecognized password hash function")
	}
	return
}

// SetPassword computes a new stored string for a password.
func SetPassword(password string) (stored string, err error) {
	n := 32768
	r := 8
	p := 1
	salt := make([]byte, 16)
	if _, err = rand.Read(salt); err != nil {
		return
	}
	var hash []byte
	if hash, err = scrypt.Key([]byte(password), salt, n, r, p, 16); err != nil {
		return
	}
	stored = fmt.Sprintf("scrypt,%d,%d,%d,%s,%s", n, r, p, base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash))
	return
}
