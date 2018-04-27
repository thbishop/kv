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
    Create(name string) error
    Delete(name string) error
    SetKey(store string, key string, value []byte) error
    GetKey(store string, key string) error
    DeleteKey(store string, key string) error
}

type etcdStore struct{
    client client.Client
    storeName string
    storePath string
}

func newEtcdStore(name string) (etcdStore, error) {
    cfg := client.Config{
        Endpoints:               []string{"http://127.0.0.1:2379"},
        Transport:               client.DefaultTransport,
        HeaderTimeoutPerRequest: time.Second,
    }
	c, err := client.New(cfg)
	if err != nil {
        return etcdStore{}, err
	}

    store := etcdStore{
        client: c,
        storeName: name,
        storePath: fmt.Sprintf("/stores/%s/", name),
    }

    return store, nil
}

func (s *etcdStore) Create() error {
	kapi := client.NewKeysAPI(s.client)
	log.Printf("Creating store '%s'\n", s.storeName)
    opts := client.SetOptions{Dir: true}
    resp, err := kapi.Set(context.Background(), s.storePath, "", &opts)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("Set is done. Metadata is %q\n", resp)
	}
    return nil
}

func (s *etcdStore) Delete() error {
	kapi := client.NewKeysAPI(s.client)
	log.Printf("Deleting store '%s'\n", s.storeName)

    opts := client.DeleteOptions{
        Dir: true,
        Recursive: true,
    }
    resp, err := kapi.Delete(context.Background(), s.storePath, &opts)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("Deletion is done. Metadata is %q\n", resp)
	}
    return nil
}

func (s *etcdStore) GetKey (name string) ([]byte, error) {
	kapi := client.NewKeysAPI(s.client)
    log.Printf("Setting key '%s' in store '%s'", name, s.storeName)
    key := fmt.Sprintf("%s/%s", s.storePath, name)
    // TODO do these make sens?
    opts := &client.GetOptions{
        Recursive: false,
        Quorum: true,
    }

    resp, err := kapi.Get(context.Background(), key, opts)
	if err != nil {
        return []byte{}, err
	}

    data, err := base64.StdEncoding.DecodeString(resp.Node.Value)
    if err != nil {
        return []byte{}, err
    }

    return data, nil
}

func (s *etcdStore) SetKey (name string, value []byte) error {
	kapi := client.NewKeysAPI(s.client)
    log.Printf("Setting key '%s' in store '%s'", name, s.storeName)
    key := fmt.Sprintf("%s/%s", s.storePath, name)
    resp, err := kapi.Set(context.Background(), key, base64.StdEncoding.EncodeToString(value), &client.SetOptions{})
	if err != nil {
        return err
	} else {
		// print common key info
		log.Printf("Set is done. Metadata is %q\n", resp)
	}
    return nil
}

func (s *etcdStore) DeleteKey (name string) error {
	kapi := client.NewKeysAPI(s.client)
    log.Printf("Deleting key '%s' in store '%s'", name, s.storeName)
    key := fmt.Sprintf("%s/%s", s.storePath, name)
    opts := &client.DeleteOptions{
        Recursive: false,
        Dir: false,
    }

    resp, err := kapi.Delete(context.Background(), key, opts)
	if err != nil {
        return err
	} else {
		// print common key info
		log.Printf("Delete is done. Metadata is %q\n", resp)
	}
    return nil
}
