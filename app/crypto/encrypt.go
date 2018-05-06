package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/jchavannes/jgo/jerr"
	"golang.org/x/crypto/scrypt"
	"io"
)

var salt = []byte{0xfe, 0xa9, 0xe9, 0x4c, 0xd9, 0x84, 0x50, 0x3d}

func SetSalt(newSalt []byte) {
	salt = newSalt
}

// https://golang.org/pkg/crypto/cipher/#example_NewCFBEncrypter
func Encrypt(value []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, jerr.Get("error getting new cipher", err)
	}

	encryptedValue := make([]byte, aes.BlockSize+len(value))
	iv := encryptedValue[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, jerr.Get("error reading random data for iv", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encryptedValue[aes.BlockSize:], value)

	return encryptedValue, nil
}

// https://golang.org/pkg/crypto/cipher/#example_NewCFBDecrypter
func Decrypt(value []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, jerr.Get("error getting new cipher", err)
	}

	if len(value) < aes.BlockSize {
		return []byte{}, jerr.New("ciphertext too short")
	}
	iv := value[:aes.BlockSize]
	value = value[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(value, value)

	return value, nil
}

// https://godoc.org/golang.org/x/crypto/scrypt#example-package
func GenerateEncryptionKeyFromPassword(password string) ([]byte, error) {
	dk, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return []byte{}, jerr.Get("error generating key", err)
	}
	return dk, nil
}
