package db

import (
	"encoding/gob"

	uuid "github.com/satori/go.uuid"
)

type EncryptionKey struct {
	Uuid         uuid.UUID
	EncryptedKey []byte
	DecryptedKey []byte
	Iterations   uint64
	Validation   []byte
}

func (key *EncryptionKey) Load(dir, fileName string) error {
	buf, err := loadFile(dir, fileName)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(buf)
	return dec.Decode(key)
}
