package paxos

import (
	"testing"
	"time"
)

func TestPaxosWithSingleProposer(t *testing.T) {
	// 1, 2, 3 are acceptors | 1001 is a proposer | 2001 is a learner
	pn := newPaxosNetwork(1, 2, 3, 1001, 2001)

	as := make([]*acceptor, 0)
	for i := 1; i < 4; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
	}

	// run accrptors
	for _, a := range as {
		go a.run()
	}

	// crate and run of proposer
	p := newProposer(1001, "hello world", pn.agentNetwork(1001), 1, 2, 3)
	go p.run()

	// create and run of learner
	l := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	value := l.learn()
	if value != "hello world" {
		t.Errorf("value = %s, want %s", value, "hello world")
	}
}

func TestPaxosWithTwoProposers(t *testing.T) {
	// 1, 2, 3 are accrptors | 1001, 1002 are proposer | 2001 is learner
	pn := newPaxosNetwork(1, 2, 3, 1001, 1002, 2001)

	as := make([]*acceptor, 0)
	for i := 1; i <= 3; i++ {
		as = append(as, newAcceptor(i, pn.agentNetwork(i), 2001))
	}

	for _, a := range as {
		go a.run()
	}

	// first proposer
	p1 := newProposer(1001, "hello world", pn.agentNetwork(1001), 1, 2, 3)
	go p1.run()

	// second proposer
	time.Sleep(time.Millisecond)
	p2 := newProposer(1002, "game over", pn.agentNetwork(1002), 1, 2, 3)
	go p2.run()

	l := newLearner(2001, pn.agentNetwork(2001), 1, 2, 3)
	value := l.learn()
	if value != "hello world" {
		t.Errorf("value = %s, want %s", value, "hello world")
	}
	time.Sleep(time.Millisecond)
}
