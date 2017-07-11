// Challenge1
package main

import (
	"fmt"
)

type PubSub struct {
	Outputs []chan string
	lock    chan struct{}
}

func NewPubSub() *PubSub {

	return &PubSub{make([]chan string, 0), make(chan struct{}, 1)}

}

func (p *PubSub) Subscribe() <-chan string {

	NewOne := make(chan string)

	p.lock <- struct{}{}
	p.Outputs = append(p.Outputs, NewOne)
	<-p.lock

	return NewOne

}

func (p *PubSub) Publish() chan<- string {

	Input := make(chan string)

	go func() {

		defer close(Input)

		msg := <-Input

		for i := range p.Outputs {
			p.Outputs[i] <- msg
		}

	}()

	return Input
}

func main() {

	ps := NewPubSub()
	a := ps.Subscribe()
	b := ps.Subscribe()
	c := ps.Subscribe()
	go func() {
		ps.Publish() <- "wat"
		ps.Publish() <- ("wat" + <-c)
	}()
	fmt.Printf("A recieved %s, B recieved %s and we ignore C!\n", <-a, <-b)
	fmt.Printf("A recieved %s, B recieved %s and C received %s\n", <-a, <-b, <-c)
}
