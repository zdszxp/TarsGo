// Package store is an interface for distribute data storage.
package store

import (
	"errors"
	"time"
	"fmt"
	"strings"

	"github.com/TarsCloud/TarsGo/tars/config/options"

	proto "github.com/golang/protobuf/proto"
)

var (
	ErrNotFound = errors.New("not found")
	ErrNotSupport = errors.New("not support")
)

func Read(s Store, x proto.Message, id string) error {
	key, err := GetKey(x, id)
	if err != nil {
		return err
	}

	record, err := s.Read(key)
	if err != nil {
		return err
	}

	if err := proto.Unmarshal(record.Value, x); err != nil {
		return err
	}

	return nil
}

func Write(s Store, x proto.Message, id string) error {
	key, err := GetKey(x, id)
	if err != nil {
		return err
	}

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

func GetKey(x proto.Message, id string) (string, error) {
	key := proto.MessageName(x)
	if len(key) <= 0 {
		return key, errors.New("GetKey Failed: check import proto path")
	}

	lastIndex := strings.LastIndex(key, ".")
	key = key[lastIndex+1:len(key)]
	if len(key) <= 0 {
		return key, errors.New(fmt.Sprintf("GetKey Failed:key %v", proto.MessageName(x)))
	}

	return key+":"+id, nil
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
