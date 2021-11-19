package paxos

import (
	"log"
	"time"
)

type network interface {
	send(m message)
	recv(timeout time.Duration) (message, bool)
}

type paxosNetwork struct {
	recvQueue map[int]chan message
}

func newPaxosNetwork(agents ...int) *paxosNetwork {
	pn := &paxosNetwork{
		recvQueue: make(map[int]chan message, 0),
	}

	for _, v := range agents {
		pn.recvQueue[v] = make(chan message, 1024)
	}
	return pn
}

func (pn *paxosNetwork) agentNetwork(id int) *agentNetwork {
	return &agentNetwork{id: id, paxosNetwork: pn}
}

func (pn *paxosNetwork) send(m message) {
	log.Printf("nt: send %+v", m)
	pn.recvQueue[m.to] <- m
}

func (pn *paxosNetwork) empty() bool {
	var n int
	for i, q := range pn.recvQueue {
		log.Printf("nt: %d left %d", i, len(q))
		n += len(q)
	}
	return n == 0
}

func (pn *paxosNetwork) recvFrom(from int, timeout time.Duration) (message, bool) {
	select {
	case m := <-pn.recvQueue[from]:
		log.Printf("nt: recv%+v", m)
		return m, true
	case <-time.After(timeout):
		return message{}, false
	}
}

type agentNetwork struct {
	id int
	*paxosNetwork
}

func (a *agentNetwork) send(m message) {
	a.paxosNetwork.send(m)
}

func (a *agentNetwork) recv(timeout time.Duration) (message, bool) {
	return a.paxosNetwork.recvFrom(a.id, timeout)
}
