package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/client"
)

type Store interface {
	CreateStore(name string) error
	DeleteStore(name string) error
	StoreExists(name string) (bool, error)
	SetKey(store string, key string, value []byte) error
	GetKey(store string, key string) ([]byte, error)
	DeleteKey(store string, key string) error
	KeyExists(store string, key string) (bool, error)
}

type etcdStore struct {
	client *client.Client
}

func newEtcdStore() (*etcdStore, error) {
	cfg := client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return &etcdStore{}, err
	}

	return &etcdStore{client: &c}, nil
}

func (s *etcdStore) storePath(storeName string) string {
	return fmt.Sprintf("/stores/%s/", storeName)
}

func (s *etcdStore) keyPath(storeName string, keyName string) string {
	return fmt.Sprintf("%s%s", s.storePath(storeName), keyName)
}

func (s *etcdStore) CreateStore(storeName string) error {
	kapi := client.NewKeysAPI(*s.client)
	log.Printf("Creating store '%s'\n", storeName)

	opts := client.SetOptions{Dir: true}
	resp, err := kapi.Set(context.Background(), s.storePath(storeName), "", &opts)
	if err != nil {
		log.Printf("Error trying to create store: %s", err)
		return err
	}

	log.Printf("Set is done. Metadata is %q\n", resp)

	return nil
}

func (s *etcdStore) DeleteStore(storeName string) error {
	kapi := client.NewKeysAPI(*s.client)
	log.Printf("Deleting store '%s'\n", storeName)

	opts := client.DeleteOptions{
		Dir:       true,
		Recursive: true,
	}
	_, err := kapi.Delete(context.Background(), s.storePath(storeName), &opts)
	if err != nil {
		log.Printf("Error trying to delete store: %s", err)
		return err
	}

	return nil
}

func (s *etcdStore) StoreExists(storeName string) (bool, error) {
	log.Printf("Checking if store exists '%s'\n", storeName)
	return s.genericKeyExists(s.storePath(storeName))
}

func (s *etcdStore) KeyExists(storeName string, keyName string) (bool, error) {
	log.Printf("Checking if key '%s' in store '%s' exists", keyName, storeName)
	return s.genericKeyExists(s.keyPath(storeName, keyName))
}

func (s *etcdStore) genericKeyExists(keyPath string) (bool, error) {
	kapi := client.NewKeysAPI(*s.client)

	// TODO do these make sens?
	opts := &client.GetOptions{
		Recursive: false,
		Quorum:    true,
	}

	_, err := kapi.Get(context.Background(), keyPath, opts)
	if err != nil {
		if client.IsKeyNotFound(err) {
			log.Printf("Generic key '%s' is not found", keyPath)
			return false, nil
		}

		log.Printf("Error trying to see if generic key exists: %s", err)
		return false, err
	}

	return true, nil
}

func (s *etcdStore) GetKey(storeName string, keyName string) ([]byte, error) {
	kapi := client.NewKeysAPI(*s.client)
	log.Printf("Getting key '%s' in store '%s'", keyName, storeName)

	// TODO do these make sens?
	opts := &client.GetOptions{
		Recursive: false,
		Quorum:    true,
	}

	resp, err := kapi.Get(context.Background(), s.keyPath(storeName, keyName), opts)
	if err != nil {
		log.Printf("Error trying to get key: %s", err)
		return []byte{}, err
	}

	data, err := base64.StdEncoding.DecodeString(resp.Node.Value)
	if err != nil {
		log.Printf("Error trying to decode key value: %s", err)
		return []byte{}, err
	}

	return data, nil
}

func (s *etcdStore) SetKey(storeName string, keyName string, value []byte) error {
	kapi := client.NewKeysAPI(*s.client)
	log.Printf("Setting key '%s' in store '%s'", keyName, storeName)

	resp, err := kapi.Set(context.Background(), s.keyPath(storeName, keyName), base64.StdEncoding.EncodeToString(value), &client.SetOptions{})
	if err != nil {
		log.Printf("Error trying to set key: %s", err)
		return err
	}

	log.Printf("Set is done. Metadata is %q\n", resp)
	return nil
}

func (s *etcdStore) DeleteKey(storeName string, keyName string) error {
	kapi := client.NewKeysAPI(*s.client)
	log.Printf("Deleting key '%s' in store '%s'", keyName, storeName)

	opts := &client.DeleteOptions{
		Recursive: false,
		Dir:       false,
	}

	resp, err := kapi.Delete(context.Background(), s.keyPath(storeName, keyName), opts)
	if err != nil {
		log.Printf("Error trying to delete key: %s", err)
		return err
	} else {
		// print common key info
		log.Printf("Delete is done. Metadata is %q\n", resp)
	}
	return nil
}
