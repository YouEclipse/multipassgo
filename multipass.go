package multipassgo

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
)

type Multipass struct {
	EncryptionKey []byte
	SignatureKey  []byte
	UserInfo      interface{}
}

func NewMultipass(secret string) *Multipass {
	//  The secret is used to derive two cryptographic keys — one for encryption and one for signing.
	//  This key derivation is done through the use of the SHA-256 hash function
	hash := sha256.New()
	hash.Write([]byte(secret))
	key := hash.Sum(nil)
	// the first 128 bit are used as encryption key
	// and the last 128 bit are used as signature key.
	return &Multipass{
		EncryptionKey: key[0:16],
		SignatureKey:  key[16:],
	}
}

func (s *Multipass) GenerateToken() (token string, err error) {
	src, err := json.Marshal(s.UserInfo)
	if err != nil {
		return
	}

	//step2:Encrypt the JSON data by AES
	ciphertext, err := encrypt([]byte(s.EncryptionKey), src)
	if err != nil {
		return
	}
	//step3:Sign the encrypted data in setp2 by HMAC
	signedtext, err := sign([]byte(s.SignatureKey), ciphertext)
	if err != nil {
		return
	}
	// The multipass login token now consists of the 128 bit initialization vector,
	// a variable length ciphertext, and a 256 bit signature (in this order).
	// This data is encoded using base64 (URL-safe variant, RFC 4648).
	t := append(ciphertext, signedtext...)
	token = base64.URLEncoding.EncodeToString([]byte(t))
	return
}

func sign(key, ciphertext []byte) (signedtext []byte, err error) {
	h := hmac.New(sha256.New, key)
	_, err = h.Write([]byte(ciphertext))
	if err != nil {
		return
	}
	signedtext = h.Sum(nil)
	return
}

func encrypt(key, text []byte) (ciphertext []byte, err error) {
	//use PKCS5Padding
	src := pkcs5Padding([]byte(text), aes.BlockSize)
	if len(src)%aes.BlockSize != 0 {
		return ciphertext, errors.New("crypto/cipher: input not full blocks")
	}

	// use the AES algorithm (128 bit key length, CBC mode of operation, random initialization vector).
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ciphertext = make([]byte, aes.BlockSize+len(src))

	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return
	}

	aesEncrypter := cipher.NewCBCEncrypter(block, iv)
	aesEncrypter.CryptBlocks(ciphertext[aes.BlockSize:], src)
	return ciphertext, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
