package minioclient

import (
	"bytes"
	"context"
	"log"
	"sync"
	"testing"

	"github.com/kiortts/mikro-kit/application"
)

func TestPutAndGet(t *testing.T) {

	log.SetFlags(log.Lshortfile)

	cfg := &Config{
		Endpoint:        "192.168.31.25:9000",
		AccessKeyID:     "minio",
		SecretAccessKey: "12345678",
		BucketName:      "bigmsgs",
	}

	mainParams := &application.MainParams{
		Ctx: context.TODO(),
		Wg:  new(sync.WaitGroup),
	}

	mc := New(cfg)
	mc.Run(mainParams)

	testData := []byte("TestData")
	key := []byte("TestKey")
	proxyData, err := mc.Put(key, testData)
	if err != nil {
		t.Error(err)
	}

	log.Printf("%s", proxyData)

	key2 := mc.GetKey(proxyData)
	data2, err := mc.Get(key2)

	log.Printf("%s", data2)

	if !bytes.Equal(testData, data2) {
		t.Error("MOCK")
	}

	err = mc.Remove(key2)
	if err != nil {
		t.Error(err)
	}

	data3, err := mc.Get(key2)
	if err == nil {
		t.Error("Empty err when specified key does not exist.")
	}
	log.Printf("%s", data3)

}
