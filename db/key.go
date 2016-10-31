package db

import uuid "github.com/satori/go.uuid"

type EncryptionKey struct {
	Uuid         uuid.UUID
	EncryptedKey []byte
	DecryptedKey []byte
	Iterations   uint64
	Validation   []byte
}
