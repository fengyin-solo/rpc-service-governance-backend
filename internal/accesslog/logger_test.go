package accesslog_test

import (
	"testing"

	"rpcgate/internal/accesslog"
	"rpcgate/internal/requestpool"
)

func TestReleasedRequestDoesNotPolluteQueuedAccessLog(t *testing.T) {
	pool := requestpool.New(1)
	logger := &accesslog.Logger{}
	first := pool.Acquire()
	first.Tenant = "tenant-a"
	first.RequestID = "request-a"
	logger.Queue(first)
	pool.Release(first)

	second := pool.Acquire()
	if second.Tenant != "" || second.RequestID != "" || len(second.Payload) != 0 {
		t.Errorf("reused request still contained tenant=%q request_id=%q", second.Tenant, second.RequestID)
	}
	second.Tenant = "tenant-b"
	second.RequestID = "request-b"
	entries := logger.Flush()
	if len(entries) != 1 || entries[0].Tenant != "tenant-a" || entries[0].RequestID != "request-a" {
		t.Errorf("queued access log changed after request reuse: %+v", entries)
	}
}
