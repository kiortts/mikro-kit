package natsclient

// import (
// 	"context"
// 	"log"
// 	"sync"
// 	"testing"
// 	"time"

// 	application "bitbucket.org/tts/digitrans/common/application"
// 	m1 "bitbucket.org/tts/digitrans/messages"
// 	"bitbucket.org/tts/digitrans/services/messageservice"
// 	minioClient "bitbucket.org/tts/digitrans/services/messageservice/minioclient"
// )

// func Test_BigMsgSend(t *testing.T) {

// 	// название и версия приложения
// 	appName := "PowerAmpApplication"
// 	appVersion := "v0.1"
// 	app := application.New(appName, appVersion)

// 	// модуль обмена большими сообщениями
// 	bigMsgServiceCfg := &minioClient.Config{
// 		Endpoint:        "192.168.31.25:9000",
// 		AccessKeyID:     "minio",
// 		SecretAccessKey: "12345678",
// 	}
// 	bigMsgService := minioClient.New(bigMsgServiceCfg) // конфигурируется переменными окружения
// 	app.Add(bigMsgService)

// 	// модуль обмена сообщениями
// 	msgServiceCfg := &Config{
// 		Url:     "",
// 		Cluster: "naTTS_cluster",
// 	}
// 	msgService := New(bigMsgService, msgServiceCfg) // конфигурируется переменными окружения
// 	app.Add(msgService)

// 	err := app.Run()
// 	if err != nil {
// 		t.Errorf(err.Error())
// 	}

// 	<-time.After(time.Second * 2)

// 	id1 := "aaa"
// 	nm1 := "bbb"
// 	subj := "testsubj"

// 	ctx := context.TODO()
// 	wg := new(sync.WaitGroup)

// 	msgHanler := func(ctx context.Context, wg *sync.WaitGroup, msg messageservice.Message) {
// 		log.Println("AAAAAAAAAAa")
// 	}
// 	msgService.HandleMessages(ctx, wg, subj, msgHanler, "")

// 	msg1 := &m1.IDList{
// 		IDs:       []string{id1},
// 		Mnemonics: []string{nm1},
// 	}

// 	err = msgService.SendMessage(msg1, subj)
// 	if err != nil {
// 		t.Errorf(err.Error())
// 	}

// 	<-time.After(time.Second * 10)

// 	t.Errorf("MOCK")

// }
