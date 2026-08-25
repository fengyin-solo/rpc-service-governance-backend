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

// Release 归还 Envelope 供下次复用。归还前必须清空身份与负载，
// 否则下一次 Acquire 的请求会读到上一个请求遗留的租户和请求号。
func (p *Pool) Release(envelope *Envelope) {
	envelope.Tenant = ""
	envelope.RequestID = ""
	envelope.Payload = nil
	select {
	case p.idle <- envelope:
	default:
	}
}
