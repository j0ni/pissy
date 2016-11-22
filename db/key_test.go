package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaddingRoundTrip(t *testing.T) {
	assert := assert.New(t)
	testBytes := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	expectedResult := []byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8, 8, 8, 8, 8, 8}

	result := pad(testBytes)
	assert.NotNil(result)
	assert.Equal(expectedResult, result)

	strippedResult, err := stripPadding(result)
	assert.Nil(err)
	assert.NotNil(strippedResult)
	assert.Equal(strippedResult, testBytes)
}

func TestPaddingIntegrityFailure(t *testing.T) {
	assert := assert.New(t)
	badPad := []byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8, 8, 8, 5, 8, 8}

	result, err := stripPadding(badPad)
	assert.Nil(result)
	assert.NotNil(err)
	assert.EqualValues("Padding integrity failure", err.Error())
}

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

	if err := key.EncryptKey(passphrase); err != nil {
		assert.Fail("EncryptKey returned an error", err.Error())
	}

	assert.NotNil(key.DecryptedKey)
	assert.NotNil(key.EncryptedKey)
	assert.NotNil(key.EncryptedKeyHMAC)

	expectedKey := key.DecryptedKey
	key.DecryptedKey = nil
	if err := key.DecryptKey(passphrase); err != nil {
		assert.Fail("DecryptKey returned an error", err.Error())
	}

	assert.NotNil(key.DecryptedKey)
	assert.NotNil(key.EncryptedKey)
	assert.NotNil(key.EncryptedKeyHMAC)
	assert.EqualValues(expectedKey, key.DecryptedKey)
}
