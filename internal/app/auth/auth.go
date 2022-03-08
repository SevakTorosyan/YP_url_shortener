package auth

import (
	"crypto/aes"
	"encoding/hex"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/utils"
)

type User struct {
	ID string `json:"id"`
}

func GenerateIdentifier() string {
	return hex.EncodeToString(utils.GenerateRandom(8))
}

func NewUser() User {
	return User{ID: GenerateIdentifier()}
}

func (u User) EncryptUserID(key []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)

	if err != nil {
		return []byte{}, nil
	}

	encrypted := make([]byte, aes.BlockSize)
	aesblock.Encrypt(encrypted, []byte(u.ID))

	return encrypted, nil
}

func (u *User) DecryptUserID(key, encrypted []byte) error {
	aesblock, err := aes.NewCipher(key)

	if err != nil {
		return err
	}

	decrypted := make([]byte, aes.BlockSize)
	aesblock.Decrypt(decrypted, encrypted)
	u.ID = string(decrypted)

	return nil
}
