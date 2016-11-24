package db

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	uuid "github.com/satori/go.uuid"
)

const EncryptionKeyFile = "encryptionKey.gob"

type Database struct {
	Path    string
	Records []Record
	Key     EncryptionKey
}

func New(path string) *Database {
	return &Database{Path: path}
}

func (db *Database) Find(uuidS string) (rec Record, err error) {
	id, err := uuid.FromString(uuidS)
	if err != nil {
		return
	}
	for _, record := range db.Records {
		if record.Uuid == id {
			return record, nil
		}
	}
	err = errors.New("UUID did not match a record")
	return
}

func (db *Database) Save() (uuids []uuid.UUID, errs []error) {
	for _, record := range db.Records {
		if err := record.Save(db.Path); err != nil {
			errs = append(errs, err)
		} else {
			uuids = append(uuids, record.Uuid)
		}
	}
	return
}

func (db *Database) Load() error {
	if found, err := Exists(db.Path); err != nil {
		return err
	} else if !found {
		return errors.New(fmt.Sprintf("path does not exist: %s", db.Path))
	}
	files, err := ioutil.ReadDir(db.Path)
	if err != nil {
		return err
	}
	recRe, err := regexp.Compile(".*\\.pissy")
	if err != nil {
		return err
	}
	for _, file := range files {
		fileName := file.Name()
		if recRe.MatchString(fileName) {
			var record Record
			err := record.Load(db.Path, fileName)
			if err != nil {
				return err
			}
			db.Records = append(db.Records, record)
		} else if fileName == "key" {
			err := db.Key.Load(db.Path, fileName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func loadFile(dir, fileName string) (*bytes.Buffer, error) {
	fullPath := fmt.Sprintf("%s/%s", dir, fileName)
	if found, err := Exists(fullPath); err != nil {
		return nil, err
	} else if !found {
		return nil, errors.New(fmt.Sprintf("file not found: %s", fullPath))
	}
	var buf bytes.Buffer
	if bs, err := ioutil.ReadFile(fullPath); err != nil {
		return nil, err
	} else if _, err := buf.Write(bs); err != nil {
		return nil, err
	}
	return &buf, nil
}

func saveFile(dir, fileName string, buf bytes.Buffer) error {
	fullPath := fmt.Sprintf("%s/%s", dir, fileName)
	return ioutil.WriteFile(fullPath, buf.Bytes(), 0600)
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
