package api

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
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
	Text     string `json:"text"`
	Font     string `json:"font"`
	Padding  string `json:"padding"`
	FontSize string `json:"fontSize"`
	Height   string `json:"height"`
}

func (g *Gen) EncodeCompressed() (string, error) {
	// JSON encode
	plaintext, err := json.Marshal(g)
	if err != nil {
		return "", err
	}

	// GZIP compress
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(plaintext); err != nil {
		return "", err
	}
	w.Close()

	// Encode to base64
	return base64.URLEncoding.EncodeToString(buf.Bytes()), nil
}

func DecodeCompressed(encoded string) (*Gen, error) {
	// Decode base64
	compressed, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	// GZIP decompress
	r, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	decompressed, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// JSON decode
	var g Gen
	if err := json.Unmarshal(decompressed, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func (g *Gen) Hex() (string, error) {
	plaintext, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(plaintext), nil
}

func GetFromHex(hexStr string) (*Gen, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	var g Gen
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, err
	}

	return &g, nil
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
