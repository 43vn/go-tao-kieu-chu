package api

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
)

var (
	key = []byte(os.Getenv("AES_KEY"))
)

type Gen struct {
	Text string `json:"text"`
	Font string `json:"font"`
}

func (gen *Gen) Encrypt() (string, error) {
	if len(key) != 32 {
		return "", errors.New("key phải có độ dài 32 byte (AES-256)")
	}

	plaintext, err := json.Marshal(gen)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	plaintext = PKCS7Padding(plaintext, aes.BlockSize)

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	// Gộp IV + ciphertext
	finalData := append(iv, ciphertext...)

	// Encode hex
	return hex.EncodeToString(finalData), nil
}

// DecryptGenHex: Giải mã từ hex string thành Gen struct
func DecryptGenHex(encryptedHex string) (Gen, error) {
	var gen Gen

	if len(key) != 32 {
		return gen, errors.New("key phải có độ dài 32 byte (AES-256)")
	}

	encrypted, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return gen, err
	}

	if len(encrypted) < aes.BlockSize {
		return gen, errors.New("dữ liệu không hợp lệ")
	}

	iv := encrypted[:aes.BlockSize]
	ciphertext := encrypted[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return gen, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext, err = PKCS7Unpadding(plaintext)
	if err != nil {
		return gen, err
	}

	err = json.Unmarshal(plaintext, &gen)
	return gen, err
}

// Padding PKCS7
func PKCS7Padding(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padding...)
}

// Unpadding PKCS7
func PKCS7Unpadding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("dữ liệu rỗng")
	}
	padLen := int(data[len(data)-1])
	if padLen > len(data) {
		return nil, errors.New("padding không hợp lệ")
	}
	return data[:len(data)-padLen], nil
}
