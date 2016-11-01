package db

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Record struct {
	Uuid        uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Notes       string
	Encrypted   []byte
	ContentHash []byte
	TypeName    string
	Dirty       bool
}

func (r Record) String() string {
	return fmt.Sprintf("<Record: %s[title=%s]>", r.Uuid, r.Title)
}

func (r Record) Save(dir string) error {
	fileName := fmt.Sprintf("%s/%s.pissy", dir, r.Uuid)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(r); err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, buf.Bytes(), 0600)
}

func (record *Record) Load(dir string, fileName string) error {
	buf, err := loadFile(dir, fmt.Sprintf("%s.pissy", fileName))
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(buf)
	return dec.Decode(record)
}

func NewRecord() Record {
	var record Record
	now := time.Now()
	record.Uuid = uuid.NewV4()
	record.CreatedAt = now
	record.UpdatedAt = now
	record.Dirty = true
	return record
}
