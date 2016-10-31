package db

import uuid "github.com/satori/go.uuid"

type Database struct {
	Path    string
	Records []Record
	Key     EncryptionKey
}

func (db *Database) Save() (uuids []uuid.UUID, errs []error) {
	for _, record := range db.Records {
		if uuid, err := record.Save(db.Path); err != nil {
			errs = append(errs, err)
		} else {
			uuids = append(uuids, uuid)
		}
	}
	return
}
