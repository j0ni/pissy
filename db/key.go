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

func (key *EncryptionKey) Save(dir, fileName string) error {
	var buf bytes.Buffer
	tmpKey := key.DecryptedKey
	key.DecryptedKey = nil
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(key); err != nil {
		return err
	}
	key.DecryptedKey = tmpKey
	return saveFile(dir, fileName, buf)
}

func (key *EncryptionKey) deriveKey(passphrase, salt []byte) []byte {
	return pbkdf2.Key(passphrase, salt, key.Iterations, AESKeySize, sha512.New)
}

func (key *EncryptionKey) Decrypt(passphrase []byte) (err error) {
	dk := key.deriveKey(passphrase, key.Salt)
	// not sure if it's ok to reuse this key
	if badHMAC(key.EncryptedKey, key.EncryptedKeyHMAC, dk) {
		return errors.New("HMAC validation failed")
	}
	key.DecryptedKey, err = decrypt(key.EncryptedKey, dk)
	return
}

func (key *EncryptionKey) Unlock(passphrase []byte) error {
	return key.Decrypt(passphrase)
}

func (key *EncryptionKey) Encrypt(passphrase []byte) error {
	key.Salt = secureRandomBytes(SaltSize)
	dk := key.deriveKey(passphrase, key.Salt)
	ciphertext, err := encrypt(key.DecryptedKey, dk)
	if err != nil {
		return err
	}
	key.EncryptedKey = ciphertext
	// not sure if it's ok to reuse this key
	key.EncryptedKeyHMAC = generateHMAC(key.EncryptedKey, dk)
	return nil
}

func (key *EncryptionKey) Lock(passphrase []byte) error {
	err := key.Encrypt(passphrase)
	if err != nil {
		return err
	}
	key.DecryptedKey = nil
	return nil
}

func (key *EncryptionKey) GenerateKey() {
	key.DecryptedKey = secureRandomBytes(AESKeySize)
}

func NewKey() *EncryptionKey {
	key := &EncryptionKey{
		Iterations: 10000,
	}
	key.GenerateKey()
	return key
}
