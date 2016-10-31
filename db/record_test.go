package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialization(t *testing.T) {
	assert := assert.New(t)

	record := NewRecord()
	record.Encrypted = []byte("yo hello there")
	record.Notes = "some notes to test"
	record.Title = "hokey kokey"
	record.Save("/tmp/wibble")

	otherRecord, err := LoadRecord("/tmp/wibble")
	assert.Nil(err, "LoadRecord returned an error")

	assert.Equal(record.String(), otherRecord.String(), "records were not equal")

	assert.Nil(os.Remove("/tmp/wibble"))
}
