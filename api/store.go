package main

import (
    "context"
    "encoding/base64"
    "log"
    "time"
    "github.com/coreos/etcd/client"
)

type Store interface {
    Create(name string) error
    SetKey(store string, key string, value []byte) error
    // Set(key string, value []byte) error
}

type etcdStore struct{
    client client.Client

}
func newEtcdStore() etcdStore {
    cfg := client.Config{
        Endpoints:               []string{"http://127.0.0.1:2379"},
        Transport:               client.DefaultTransport,
        // set timeout per request to fail fast when the target endpoint is unavailable
        HeaderTimeoutPerRequest: time.Second,
    }
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
    return etcdStore{client: c}
}

func (s *etcdStore) Create(name string) error {
	kapi := client.NewKeysAPI(s.client)
    opts := client.SetOptions{Dir: true}
	// set "/foo" key with "bar" value
	log.Printf("Creating store '%s'\n", name)
    resp, err := kapi.Set(context.Background(), name, "", &opts)
	// resp, err := kapi.Set(context.Background(), key, "bar", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("Set is done. Metadata is %q\n", resp)
	}
    return nil
}

func (s *etcdStore) SetKey (store string, name string, value []byte) error {
	kapi := client.NewKeysAPI(s.client)
    log.Printf("Setting key '%s' in store '%s'", name, store)
    key := "/" + store + "/" + name
    resp, err := kapi.Set(context.Background(), key, base64.StdEncoding.EncodeToString(value), &client.SetOptions{})
	if err != nil {
        return err
	} else {
		// print common key info
		log.Printf("Set is done. Metadata is %q\n", resp)
	}
    return nil
}
