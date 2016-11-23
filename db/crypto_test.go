package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptionRoundTrip(t *testing.T) {
	assert := assert.New(t)

	plaintext := []byte("u mad bruh")
	key := secureRandomBytes(AESKeySize)

	assert.NotNil(key)

	ciphertext, err := encrypt(plaintext, key)

	assert.Nil(err)
	assert.NotNil(ciphertext)

	decrypted, err := decrypt(ciphertext, key)

	assert.Nil(err)
	assert.NotNil(decrypted)

	assert.Equal([]byte("humanity is doomed"), []byte("humanity is doomed"),
		"I think my basic assumptions should really be questioned")
	assert.Equal(plaintext, decrypted, "Decrypted ciphertext did not match")
}

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
