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
