package helloworld

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/kiortts/mikro-kit/application"
)

func TestPrintHelloFunc(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	name := "Test"
	wg := new(sync.WaitGroup)
	wg.Add(1)
	printHello(context.TODO(), wg, name)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	ref := fmt.Sprintf("Hello, %s!!!\n", name)
	if string(out) != ref {
		t.Errorf("Wrong out text: %s", out)
	}
}

func TestHelloWorld(t *testing.T) {

	name := "Test"
	cfg := &Config{
		Name: name,
	}
	helloWorldInstance := New(cfg) // current configuration

	CommonHelloWorld(t, name, helloWorldInstance)
}

func TestDefaultConfig(t *testing.T) {

	name := "World"                // default value
	helloWorldInstance := New(nil) // self configuration

	CommonHelloWorld(t, name, helloWorldInstance)
}

func TestEnvConfig(t *testing.T) {

	name := "Test"
	os.Setenv("HELLO_NAME", name)  // set name to env
	helloWorldInstance := New(nil) // self configuration

	CommonHelloWorld(t, name, helloWorldInstance)
}

func CommonHelloWorld(t *testing.T, name string, instance *HelloWordModule) {

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	mainParams := &application.MainParams{
		Ctx:     context.TODO(),
		Wg:      new(sync.WaitGroup),
		AppStop: func() {},
	}
	err := instance.Run(mainParams)
	if err != nil {
		t.Errorf("HelloWorld Run err: %s", err)
	}
	<-time.After(time.Millisecond * 100)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	ref := fmt.Sprintf("Hello, %s!!!\n", name)
	if string(out) != ref {
		t.Errorf("Wrong out text: %s", out)
	}
}
