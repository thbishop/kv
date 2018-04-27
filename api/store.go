package main

import (
    "context"
    "log"
    "time"
    "github.com/coreos/etcd/client"
)

type Store interface {
    Create(name string) error
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
	log.Print("Setting '/foo' key with 'bar' value")
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
