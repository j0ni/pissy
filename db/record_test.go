package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialization(t *testing.T) {
	assert := assert.New(t)

	record := NewRecord()
	record.EncryptedValue = []byte("yo hello there")
	record.Notes = "some notes to test"
	record.Title = "hokey kokey"
	record.Save("/tmp")

	var otherRecord Record
	err := otherRecord.Load("/tmp", record.Uuid.String())
	assert.Nil(err, "LoadRecord returned an error")

	assert.Equal(record.String(), otherRecord.String(), "records were not equal")

	assert.Nil(os.Remove(fmt.Sprintf("/tmp/%s.pissy", record.Uuid.String())))
}

func TestRoundTripRecordEncryption(t *testing.T) {
	assert := assert.New(t)
	plaintext := []byte("u mad bruh")

	var key EncryptionKey
	key.GenerateKey()
	assert.NotNil(key.DecryptedKey)

	record := NewRecord()
	record.DecryptedValue = plaintext

	assert.Nil(record.EncryptedValue)
	assert.Nil(record.EncryptedValueHMAC)

	err := record.Encrypt(key.DecryptedKey)

	assert.Nil(err)
	assert.NotNil(record.EncryptedValue)
	assert.NotNil(record.EncryptedValueHMAC)

	record.DecryptedValue = nil
	err = record.Decrypt(key.DecryptedKey)

	assert.Nil(err)
	assert.NotNil(record.DecryptedValue)
	assert.NotNil(record.EncryptedValue)
	assert.NotNil(record.EncryptedValueHMAC)
	assert.Equal(plaintext, record.DecryptedValue)
}
