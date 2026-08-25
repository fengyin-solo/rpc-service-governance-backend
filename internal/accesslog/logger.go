package accesslog

import "rpcgate/internal/requestpool"

type Entry struct {
	Tenant    string
	RequestID string
}

// Logger 收集访问日志。Queue 时立即对 Envelope 的身份字段做快照，
// 这样即使该 Envelope 后续被对象池复用、字段被改写或清空，
// 已排队的日志条目身份也不会再变。
type Logger struct{ queued []Entry }

func (l *Logger) Queue(envelope *requestpool.Envelope) {
	l.queued = append(l.queued, Entry{
		Tenant:    envelope.Tenant,
		RequestID: envelope.RequestID,
	})
}

func (l *Logger) Flush() []Entry {
	entries := make([]Entry, 0, len(l.queued))
	entries = append(entries, l.queued...)
	l.queued = l.queued[:0]
	return entries
}
