// Package store is an interface for distribute data storage.
package store

import (
	"errors"
	"time"

	"github.com/TarsCloud/TarsGo/tars/config/options"

	proto "github.com/golang/protobuf/proto"
)

var (
	ErrNotFound = errors.New("not found")
	ErrNotSupport = errors.New("not support")
)

func Read(s Store, x proto.Message, id string) error {
	key := GetKey(x, id)

	record, error := s.Read(key)
	if error != nil {
		return error
	}

	if err := proto.Unmarshal(record.Value, x); err != nil {
		return err
	}

	return nil
}

func Write(s Store, x proto.Message, id string) error {
	key := GetKey(x, id)

	data, err := proto.Marshal(x)
	if err != nil {
		return err
	}
	
	record := &Record{
		Key : key,
		Value : data,
	}

	error := s.Write(record)
	if error != nil {
		return error
	}

	return nil
}

func GetKey(x proto.Message, id string) string {
	return proto.MessageName(x)+":"+id
}

// Store is a data storage interface
type Store interface {
	// embed options
	options.Options
	// Dump the known records
	Dump() ([]*Record, error)
	// Read a record with key
	Read(key string) (*Record, error)
	// Write a record
	Write(r *Record) error
	// Delete a record with key
	Delete(key string) error
}

// Record represents a data record
type Record struct {
	Key    string
	Value  []byte
	Expiry time.Duration
}
