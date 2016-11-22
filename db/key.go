package db

import (
	"bytes"
	"crypto/sha512"
	"encoding/gob"
	"errors"

	"golang.org/x/crypto/pbkdf2"

	uuid "github.com/satori/go.uuid"
)

type EncryptionKey struct {
	Uuid             uuid.UUID
	EncryptedKey     []byte
	DecryptedKey     []byte
	EncryptedKeyHMAC []byte
	Salt             []byte
	Iterations       int
}

func (key *EncryptionKey) Load(dir, fileName string) error {
	buf, err := loadFile(dir, fileName)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(buf)
	return dec.Decode(key)
}

func (key EncryptionKey) Save(dir, fileName string) error {
	var buf bytes.Buffer
	// its a copy
	key.DecryptedKey = nil
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(key); err != nil {
		return err
	}
	return saveFile(dir, fileName, buf)
}

func (key EncryptionKey) deriveKey(passphrase string, salt []byte) []byte {
	return pbkdf2.Key([]byte(passphrase), salt, key.Iterations, 32, sha512.New)
}

func (key *EncryptionKey) DecryptKey(passphrase string) error {
	dk := key.deriveKey(passphrase, key.Salt)
	// not sure if it's ok to reuse this key
	if !checkHMAC(key.EncryptedKey, key.EncryptedKeyHMAC, dk) {
		return errors.New("HMAC validation failed")
	}
	plaintext, err := decrypt(dk, key.EncryptedKey)
	key.DecryptedKey = plaintext
	return err
}

func (key *EncryptionKey) EncryptKey(passphrase string) error {
	key.Salt = secureRandomBytes(SaltSize)
	dk := key.deriveKey(passphrase, key.Salt)
	ciphertext, err := encrypt(dk, key.DecryptedKey)
	if err != nil {
		return err
	}
	key.EncryptedKey = ciphertext
	// not sure if it's ok to reuse this key
	key.EncryptedKeyHMAC = generateHMAC(key.EncryptedKey, dk)
	return nil
}

func (key *EncryptionKey) GenerateKey() {
	key.DecryptedKey = secureRandomBytes(AESKeySize)
}
