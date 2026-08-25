package requestpool

type Envelope struct {
	Tenant    string
	RequestID string
	Payload   []byte
}

type Pool struct{ idle chan *Envelope }

func New(size int) *Pool { return &Pool{idle: make(chan *Envelope, size)} }

func (p *Pool) Acquire() *Envelope {
	select {
	case envelope := <-p.idle:
		return envelope
	default:
		return &Envelope{}
	}
}

func (p *Pool) Release(envelope *Envelope) {
	select {
	case p.idle <- envelope:
	default:
	}
}
