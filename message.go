package paxos

type msgType int

const (
	Prepare msgType = iota + 1 // Send from proposer -> acceptor
	Promise                    // Send from acceptor -> proposer
	Propose                    // Send from proposer -> acceptor
	Accept                     // Send from proposer -> learner
)

type message struct {
	from  int
	to    int
	typ   msgType
	num   int
	pre   int
	value string
}

type promise interface {
	number() int
}

type accept interface {
	proposalValue() string
	proposalNumber() int
}

func (m message) number() int {
	return m.num
}

func (m message) proposalValue() string {
	switch m.typ {
	case Promise, Accept:
		return m.value
	default:
		panic("unexpected proposal value")
	}
}

func (m message) proposalNumber() int {
	switch m.typ {
	case Promise:
		return m.pre
	case Accept:
		return m.num
	default:
		panic("unexpected proposal value")
	}
}
