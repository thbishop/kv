package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

type Store interface {
	CreateStore(name string) error
	DeleteStore(name string) error
	StoreExists(name string) (bool, error)
	SetKey(store string, key string, value []byte) error
	GetKey(store string, key string) ([]byte, error)
	DeleteKey(store string, key string) error
	KeyExists(store string, key string) (bool, error)
	IsKeyMissing(error) bool
}

type consulStore struct {
	kv *api.KV
}

func newConsulStore() (*consulStore, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return &consulStore{}, err
	}

	return &consulStore{kv: client.KV()}, nil
}

func (s *consulStore) StoreExists(storeName string) (bool, error) {
	log.Printf("Checking if store exists '%s'", storeName)
	return s.genericKeyExists(s.storePath(storeName))
}

func (s *consulStore) KeyExists(storeName string, keyName string) (bool, error) {
	log.Printf("Checking if key '%s' in store '%s' exists", keyName, storeName)
	return s.genericKeyExists(s.keyPath(storeName, keyName))
}

func (s *consulStore) IsKeyMissing(err error) bool {
	return err.Error() == "key not found"
}

func (s *consulStore) CreateStore(storeName string) error {
	log.Printf("Creating store '%s'\n", storeName)
	_, err := s.kv.Put(&api.KVPair{Key: s.storePath(storeName)}, nil)
	if err != nil {
		log.Printf("Error trying to create store: %s", err)
		return err
	}

	log.Printf("Store '%s' created", storeName)
	return nil
}

func (s *consulStore) DeleteStore(storeName string) error {
	log.Printf("Deleting store '%s'\n", storeName)

	_, err := s.kv.DeleteTree(s.storePath(storeName), nil)
	if err != nil {
		log.Printf("Error trying to delete store: %s", err)
		return err
	}

	log.Printf("Store '%s' deleted", storeName)
	return nil
}

func (s *consulStore) SetKey(storeName string, keyName string, value []byte) error {
	log.Printf("Setting key '%s' in store '%s'", keyName, storeName)

	key := s.keyPath(storeName, keyName)

	_, err := s.kv.Put(&api.KVPair{Key: key, Value: value}, nil)
	if err != nil {
		log.Printf("Error trying to set key: %s", err)
		return err
	}

	log.Printf("Set key '%s'", key)
	return nil
}

func (s *consulStore) GetKey(storeName string, keyName string) ([]byte, error) {
	log.Printf("Getting key '%s' in store '%s'", keyName, storeName)

	pair, _, err := s.kv.Get(s.keyPath(storeName, keyName), nil)
	if err != nil {
		log.Printf("Error trying to get key: %s", err)
		return []byte{}, err
	}

	if pair == nil {
		log.Printf("Key '%s' not found", keyName)
		return []byte{}, errors.New("key not found")
	}

	log.Printf("Retrieved key '%s'", keyName)
	return pair.Value, nil
}

func (s *consulStore) DeleteKey(storeName string, keyName string) error {
	log.Printf("Deleting key '%s' in store '%s'", keyName, storeName)

	key := s.keyPath(storeName, keyName)
	_, err := s.kv.Delete(key, nil)
	if err != nil {
		log.Printf("Error trying to delete key: %s", err)
		return err
	}

	log.Printf("Deleted key '%s'", key)
	return nil
}

func (s *consulStore) storePath(storeName string) string {
	return fmt.Sprintf("stores/%s", storeName)
}

func (s *consulStore) keyPath(storeName string, keyName string) string {
	return fmt.Sprintf("%s/%s", s.storePath(storeName), keyName)
}

func (s *consulStore) genericKeyExists(keyPath string) (bool, error) {
	pair, _, err := s.kv.Get(keyPath, nil)
	if err != nil {
		log.Printf("Error trying to see if generic key '%s' exists: %s", keyPath, err)
		return false, err
	}

	return pair != nil, nil
}
