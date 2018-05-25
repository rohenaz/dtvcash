package crypto_test

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/rohenaz/dtvcash/app/crypto"
	"testing"
)

const Password = "bl&8aZ&74CK!vMBt"
const SomePlaintext = "some plaintext"

func TestEncryption(t *testing.T) {
	key, err := crypto.GenerateEncryptionKeyFromPassword(Password)
	if err != nil {
		t.Fatal(jerr.Get("error generating key from password", err))
	}
	fmt.Printf("Key: %x\n", key)
	encrypted, err := crypto.Encrypt([]byte(SomePlaintext), key)
	if err != nil {
		t.Fatal(jerr.Get("failed to encrypt", err))
	}
	fmt.Printf("Encrypted: %x\n", encrypted)
	decrypted, err := crypto.Decrypt(encrypted, key)
	if err != nil {
		t.Fatal(jerr.Get("failed to decrypt", err))
	}
	fmt.Printf("Decrypted: %s\n", decrypted)
	if string(decrypted) != SomePlaintext {
		t.Fatal(jerr.New("decrypted value does not match original"))
	}
}
