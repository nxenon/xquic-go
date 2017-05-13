package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"time"

	"golang.org/x/crypto/hkdf"
)

// StkSource is used to create and verify source address tokens
type StkSource interface {
	// NewToken creates a new token for a given IP address
	NewToken(sourceAddress []byte) ([]byte, error)
	// DecodeToken decodes a token and returns the source address and the timestamp
	DecodeToken(data []byte) ([]byte, time.Time, error)
}

type sourceAddressToken struct {
	data []byte
	// unix timestamp in seconds
	timestamp uint64
}

func (t *sourceAddressToken) serialize() []byte {
	res := make([]byte, 8+len(t.data))
	binary.LittleEndian.PutUint64(res, t.timestamp)
	copy(res[8:], t.data)
	return res
}

func parseToken(data []byte) (*sourceAddressToken, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("STK too short: %d", len(data))
	}
	return &sourceAddressToken{
		data:      data[8:],
		timestamp: binary.LittleEndian.Uint64(data),
	}, nil
}

type stkSource struct {
	aead cipher.AEAD
}

const stkKeySize = 16

// Chrome currently sets this to 12, but discusses changing it to 16. We start
// at 16 :)
const stkNonceSize = 16

// NewStkSource creates a source for source address tokens
func NewStkSource() (StkSource, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, err
	}
	key, err := deriveKey(secret)
	if err != nil {
		return nil, err
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCMWithNonceSize(c, stkNonceSize)
	if err != nil {
		return nil, err
	}
	return &stkSource{aead: aead}, nil
}

func (s *stkSource) NewToken(data []byte) ([]byte, error) {
	return encryptToken(s.aead, &sourceAddressToken{
		data:      data,
		timestamp: uint64(time.Now().Unix()),
	})
}

func (s *stkSource) DecodeToken(data []byte) ([]byte, time.Time, error) {
	if len(data) < stkNonceSize {
		return nil, time.Time{}, errors.New("STK too short")
	}
	nonce := data[:stkNonceSize]
	res, err := s.aead.Open(nil, nonce, data[stkNonceSize:], nil)
	if err != nil {
		return nil, time.Time{}, err
	}

	token, err := parseToken(res)
	if err != nil {
		return nil, time.Time{}, err
	}
	if token.timestamp > math.MaxInt64 {
		return nil, time.Time{}, errors.New("invalid timestamp")
	}
	return token.data, time.Unix(int64(token.timestamp), 0), nil
}

func deriveKey(secret []byte) ([]byte, error) {
	r := hkdf.New(sha256.New, secret, nil, []byte("QUIC source address token key"))
	key := make([]byte, stkKeySize)
	if _, err := io.ReadFull(r, key); err != nil {
		return nil, err
	}
	return key, nil
}

func encryptToken(aead cipher.AEAD, token *sourceAddressToken) ([]byte, error) {
	nonce := make([]byte, stkNonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return aead.Seal(nonce, nonce, token.serialize(), nil), nil
}
