package cmd

import (
	"errors"
	"fmt"
	"github.com/99designs/keyring"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/argon2"
)

var (
	username string
	password string
)
var KeyGenCmd = &cobra.Command{
	Use:   "key-gen",
	Short: "Generate a new key",
	Long:  "Generate a new key to encrypt and decrypt it",
	Run: func(cmd *cobra.Command, args []string) {
		KeyGen()
	},
}

func KeyGen() {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "key-gen",
	})
	if err != nil {
		fmt.Println("Error opening keyring:", err)
		return
	}

	// Check if key exists
	_, err = ring.Get("cipher-key")
	if errors.Is(err, keyring.ErrKeyNotFound) {
		fmt.Println("Generating new encryption key...")
		key := keygen(username, password)

		err := ring.Set(keyring.Item{
			Key:         "cipher-key",
			Data:        key,
			Label:       "File encryption key",
			Description: "Key used for encrypting/decrypting files",
		})
		if err != nil {
			fmt.Println("Error saving key:", err)
			return
		}
		fmt.Println("New encryption key generated and stored securely")
	} else if err != nil {
		fmt.Println("Error checking keyring:", err)
		return
	} else {
		fmt.Println("Encryption key already exists in keyring")
	}
}
func keygen(username, password string) []byte {
	key := argon2.Key([]byte(username), []byte(password), 3, 32*1024, 4, 32)
	return key
}
