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
	if exists, err := exists(fileName); err != nil {
		return err
	} else if !exists || r.Dirty {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(r); err != nil {
			return err
		}
		return ioutil.WriteFile(fileName, buf.Bytes(), 0600)
	}
	return nil
}

func (r *Record) Load(dir string, fileName string) error {
	buf, err := loadFile(dir, fmt.Sprintf("%s.pissy", fileName))
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
		Dirty:     true,
	}
}
