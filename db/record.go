package db

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Record struct {
	Uuid               uuid.UUID
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Title              string
	Notes              string
	EncryptedValue     []byte
	DecryptedValue     []byte
	EncryptedValueHMAC []byte
	TypeName           string
}

func (r Record) String() string {
	return fmt.Sprintf("<Record: %s[title=%s]>", r.Uuid, r.Title)
}

func (r Record) Save(dir string) error {
	fileName := fmt.Sprintf("%s/%s.pissy", dir, r.Uuid)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	r.DecryptedValue = nil
	if err := enc.Encode(r); err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, buf.Bytes(), 0600)
}

func (r *Record) Load(dir string, fileName string) error {
	buf, err := loadFile(dir, fileName)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(buf)
	return dec.Decode(r)
}

func NewRecord() Record {
	now := time.Now()
	uuid := uuid.NewV4()
	return Record{
		Uuid:      uuid,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (r *Record) Decrypt(key []byte) error {
	if badHMAC(r.EncryptedValue, r.EncryptedValueHMAC, key) {
		return errors.New("HMAC validation failed")
	}
	plaintext, err := decrypt(r.EncryptedValue, key)
	r.DecryptedValue = plaintext
	return err
}

func (r *Record) Encrypt(key []byte) error {
	ciphertext, err := encrypt(r.DecryptedValue, key)
	if err != nil {
		return err
	}
	r.EncryptedValue = ciphertext
	// key reuse again - ok, or not ok?
	r.EncryptedValueHMAC = generateHMAC(r.EncryptedValue, key)
	return nil
}

func (r *Record) UpdateFields(key []byte, rec *Record) {
	if len(rec.Title) != 0 {
		r.Title = rec.Title
	}
	if len(rec.TypeName) != 0 {
		r.TypeName = rec.TypeName
	}
	if len(rec.Notes) != 0 {
		r.Notes = rec.Notes
	}
	if len(rec.DecryptedValue) != 0 {
		r.DecryptedValue = rec.DecryptedValue
		r.Encrypt(key)
	}
}
