package accesslog

import "rpcgate/internal/requestpool"

type Entry struct {
	Tenant    string
	RequestID string
}

type Logger struct{ queued []*requestpool.Envelope }

func (l *Logger) Queue(envelope *requestpool.Envelope) {
	l.queued = append(l.queued, envelope)
}

func (l *Logger) Flush() []Entry {
	entries := make([]Entry, 0, len(l.queued))
	for _, envelope := range l.queued {
		entries = append(entries, Entry{Tenant: envelope.Tenant, RequestID: envelope.RequestID})
	}
	return entries
}
