package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyGeneration(t *testing.T) {
	assert := assert.New(t)

	var key EncryptionKey
	assert.Nil(key.DecryptedKey)
	key.GenerateKey()
	assert.NotNil(key.DecryptedKey)

	assert.Len(key.DecryptedKey, 32)
}

func TestKeyEncryption(t *testing.T) {
	assert := assert.New(t)
	var key EncryptionKey
	passphrase := "humanity is doomed"

	assert.Nil(key.EncryptedKey)
	assert.Nil(key.DecryptedKey)
	assert.Nil(key.EncryptedKeyHMAC)

	key.GenerateKey()

	assert.NotNil(key.DecryptedKey)
	assert.Len(key.DecryptedKey, 32)
	assert.Nil(key.EncryptedKey)
	assert.Nil(key.EncryptedKeyHMAC)

	if err := key.Encrypt(passphrase); err != nil {
		assert.Fail("EncryptKey returned an error", err.Error())
	}

	assert.NotNil(key.DecryptedKey)
	assert.NotNil(key.EncryptedKey)
	assert.NotNil(key.EncryptedKeyHMAC)

	expectedKey := key.DecryptedKey
	key.DecryptedKey = nil
	if err := key.Decrypt(passphrase); err != nil {
		assert.Fail("DecryptKey returned an error", err.Error())
	}

	assert.NotNil(key.DecryptedKey)
	assert.NotNil(key.EncryptedKey)
	assert.NotNil(key.EncryptedKeyHMAC)
	assert.EqualValues(expectedKey, key.DecryptedKey)
}
