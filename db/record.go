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
}

func (r Record) String() string {
	return fmt.Sprintf("<Record: %s[title=%s]>", r.Uuid, r.Title)
}

func (r Record) Save(dir string) (uuid uuid.UUID, err error) {
	uuid = r.Uuid
	fileName := fmt.Sprintf("%s/%s", dir, uuid)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err = enc.Encode(r); err != nil {
		return
	}
	err = ioutil.WriteFile(fileName, buf.Bytes(), 0600)
	return
}

func LoadRecord(dir string, uuid uuid.UUID) (Record, error) {
	var record Record
	var buf bytes.Buffer
	fileName := fmt.Sprintf("%s/%s.gob", dir, uuid)
	if bs, err := ioutil.ReadFile(fileName); err != nil {
		return record, err
		// } else if ... crypto {
	} else if _, err := buf.Write(bs); err != nil {
		return record, err
	} else {
		dec := gob.NewDecoder(&buf)
		err = dec.Decode(&record)
	}
	return record, nil
}

func NewRecord() Record {
	var record Record
	now := time.Now()
	record.Uuid = uuid.NewV4()
	record.CreatedAt = now
	record.UpdatedAt = now
	return record
}
