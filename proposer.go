package paxos

import (
	"log"
	"time"
)

type proposer struct {
	id int
	// stable
	lastSeq      int
	proposeNum   int
	proposeValue string
	acceptors    map[int]promise
	nt           network
}

func newProposer(id int, value string, nt network, acceptors ...int) *proposer {
	p := &proposer{id: id, nt: nt, lastSeq: 0, proposeValue: value, acceptors: make(map[int]promise)}
	for _, a := range acceptors {
		p.acceptors[a] = message{}
	}
	return p
}

func (p *proposer) run() {
	var ok bool
	var m message

	// Stage one: do prepare until reach the majority
	for !p.majorityReached() {
		if !ok {
			ms := p.prepare()
			for i := range ms {
				p.nt.send(ms[i])
			}
		}
		m, ok = p.nt.recv(time.Second)
		if !ok {
			// the previoud prepare is failed
			// continue to do another prepare
			continue
		}
		switch m.typ {
		case Promise:
			p.receivePromise(m)
		default:
			log.Panicf("proposer: %d unexpected message type: %v", p.id, m.typ)
		}
	}
	log.Printf("proposer: %d promise %d reached majority %d", p.id, p.n(), p.majority())

	// Stage two: do propose
	log.Printf("proposer: %d starts to propose [%d: %s]", p.id, p.n(), p.proposeValue)
	ms := p.propose()
	for i := range ms {
		p.nt.send(ms[i])
	}
}

// If the proposer receives the requested responses from a majority of
// the acceptors, then it can issue a proposal with number n and value
// v, where v is the value of the highest-numbered proposal among the
// responses, or is any value selected by proposer if the responders
// reported no proposals.
func (p *proposer) propose() []message {
	ms := make([]message, p.majority())

	i := 0
	for to, promise := range p.acceptors {
		if promise.number() == p.n() {
			ms[i] = message{from: p.id, to: to, typ: Propose, num: p.n(), value: p.proposeValue}
			i++
		}
		if i == p.majority() {
			break
		}
	}
	return ms
}

// A proposer chooses a new proposal number n and sends a request to
// each number of some set of acceptors, asking it to respond with:
// (a) A promise never again to accept to proposal numbered less than n, and
// (b) The proposal with the highest number less than n that it has accepted, if any.
func (p *proposer) prepare() []message {
	p.lastSeq++

	ms := make([]message, p.majority())
	i := 0
	for to := range p.acceptors {
		ms[i] = message{from: p.id, to: to, typ: Prepare, num: p.n()}
		i++
		if i == p.majority() {
			break
		}
	}
	return ms
}

func (p *proposer) majorityReached() bool {
	m := 0
	for _, promise := range p.acceptors {
		if promise.number() == p.n() {
			m++
		}
	}
	return m >= p.majority()
}

func (p *proposer) n() int {
	return p.lastSeq<<16 | p.id
}

func (p *proposer) majority() int {
	return len(p.acceptors)/2 + 1
}

func (p *proposer) receivePromise(promise message) {
	prevPromise := p.acceptors[promise.from]

	if prevPromise.number() < promise.number() {
		log.Printf("proposer: %d received a new promise %+v", p.id, promise)
		p.acceptors[promise.from] = promise

		// update value to the value with a larger Num
		if promise.proposalNumber() > p.proposeNum {
			p.proposeNum = promise.proposalNumber()
			p.proposeValue = promise.proposalValue()
		}
	}
}
