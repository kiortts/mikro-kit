package natsclient_test

import (
	"bytes"
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/components/messages"
	"github.com/kiortts/mikro-kit/components/messages/messageservice/natsclient"
)

func beforeEach() (context.Context, *sync.WaitGroup, *natsclient.NatsClient) {

	log.SetFlags(log.Lshortfile)

	cfg := &natsclient.Config{
		Url:              "",
		StreamingCluster: "naTTS_cluster",
	}
	ctx := context.TODO()
	wg := new(sync.WaitGroup)
	mainParams := &application.MainParams{
		Ctx:     ctx,
		Wg:      wg,
		AppStop: func() {},
	}
	natsClient := natsclient.New(cfg, nil)
	natsClient.Run(mainParams)
	<-time.After(time.Millisecond * 100)

	return ctx, wg, natsClient
}

func Test_Send_Handle(t *testing.T) {

	// get connection to Nats server
	ctx, wg, natsClient := beforeEach()

	subj := "TestSubject"
	dataText := "TestDataText"
	var data1, data2 []byte
	data1 = []byte(dataText)
	handleFunc := func(ctx context.Context, wg *sync.WaitGroup, msg messages.Message) {
		data2 = msg.Data()
	}
	natsClient.HandleMessages(ctx, wg, subj, handleFunc)

	<-time.After(time.Millisecond * 100)

	err := natsClient.SendMessage(subj, data1)
	if err != nil {
		t.Errorf("SendMessage err %s", err)
	}

	<-time.After(time.Millisecond * 100)

	if !bytes.Equal(data1, data2) {
		t.Errorf(`Data1 not equal data2 %s -> %s`, data1, data2)
	}

}

func Test_StreamingSend_StreamingHandle(t *testing.T) {

	// get connection to Nats server
	ctx, wg, natsClient := beforeEach()

	// get subscription for subject

	subj := "TestSubject"
	var data2 []byte
	handleFunc := func(ctx context.Context, wg *sync.WaitGroup, msg messages.Message) {
		data2 = msg.Data()
	}
	natsClient.HandleStreamingMessages(ctx, wg, subj, handleFunc)
	<-time.After(time.Millisecond * 100)

	// send test message

	dataText := "TestDataText"
	data1 := []byte(dataText)
	err := natsClient.SendStreamingMessage(subj, data1)
	if err != nil {
		t.Errorf("SendMessage err %s", err)
	}
	<-time.After(time.Millisecond * 100)

	// check

	if !bytes.Equal(data1, data2) {
		t.Errorf(`Data1 not equal data2 %s -> %s`, data1, data2)
	}

}

func Test_Request_Handle(t *testing.T) {

	// get connection to Nats server
	ctx, wg, natsClient := beforeEach()

	// get subscription for subject

	subj := "TestSubject"
	respText := []byte("TestRespText")
	var data2 []byte
	handleFunc := func(ctx context.Context, wg *sync.WaitGroup, msg messages.Message) {
		data2 = append(msg.Data(), respText...)
		msg.Respond(data2)
	}
	natsClient.HandleMessages(ctx, wg, subj, handleFunc)
	<-time.After(time.Millisecond * 100)

	// send test message

	dataText := "TestDataText"
	data1 := []byte(dataText)
	resp, err := natsClient.SendRequest(subj, data1, time.Second*3)
	if err != nil {
		t.Errorf("SendSendRequest err %s", err)
	}
	<-time.After(time.Millisecond * 100)

	// check

	log.Printf("%s", data2)

	if !bytes.Equal(data2, resp.Data()) {
		t.Errorf(`Data2 not equal Respdata %s -> %s`, data1, data2)
	}

}

func Test_Send_Sub(t *testing.T) {

	// get connection to Nats server
	_, _, natsClient := beforeEach()

	subj := "TestSubject"
	dataText := "TestDataText"
	var data1 []byte
	data1 = []byte(dataText)

	in, _ := natsClient.Sub(subj)
	<-time.After(time.Millisecond * 100)

	err := natsClient.SendMessage(subj, data1)
	if err != nil {
		t.Errorf("SendMessage err %s", err)
	}

	select {
	case <-time.After(time.Second):
		t.Errorf(`Timeout`)
	case msg := <-in:
		if !bytes.Equal(data1, msg.Data()) {
			log.Println(data1)
			log.Println(msg.Data())
			t.Errorf(`data1 not equal data2: %s -> %s`, data1, msg.Data())
		}
	}
}

func Test_StreamingSend_StreamingSub(t *testing.T) {

	// get connection to Nats server
	_, _, natsClient := beforeEach()

	subj := "TestSubject"
	dataText := "TestDataText"
	var data1 []byte
	data1 = []byte(dataText)

	in, _ := natsClient.SubStreaming(subj)
	<-time.After(time.Millisecond * 100)

	err := natsClient.SendStreamingMessage(subj, data1)
	if err != nil {
		t.Errorf("SendMessage err %s", err)
	}

	select {
	case <-time.After(time.Second):
		t.Errorf(`Timeout`)
	case msg := <-in:
		if !bytes.Equal(data1, msg.Data()) {
			log.Println(data1)
			log.Println(msg.Data())
			t.Errorf(`data1 not equal data2: %s -> %s`, data1, msg.Data())
		}
	}
}
