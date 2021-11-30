package main

import (
	"fmt"
	"github.com/BGrewell/transit/structure"
	"net"
	"time"
)

func main() {

	type TestStruct struct {
		Name string `json:"name"`
		Value string `json:"value"`
	}

	done := make(chan bool)

	// This simple test application transmits a structure then receives it
	go func(done chan bool) () {
		listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 21345})
		if err != nil {
            panic(err)
        }
		conn, err := listener.Accept()
		if err != nil {
            panic(err)
        }

		var ts TestStruct
		err = structure.Receive(&ts, conn)
		if err != nil {
            panic(err)
        }

		fmt.Printf("name: %s value: %s\n", ts.Name, ts.Value)
		done <- true
		listener.Close()
	}(done)
	time.Sleep(1 * time.Second)

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: 21345})
	if err != nil {
        panic(err)
    }

	ts := TestStruct{Name: "test", Value: "some value"}
	err = structure.Transmit(&ts, conn)
	if err != nil {
        panic(err)
    }

	<-done
}
