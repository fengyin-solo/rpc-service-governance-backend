package publish

import "rpcgate/internal/session"

type Runner func(string) error

type Batch struct {
	manager *session.Manager
	runner  Runner
	success []string
}

func New(manager *session.Manager, runner Runner) *Batch {
	return &Batch{manager: manager, runner: runner}
}

func (b *Batch) Process(names []string) error {
	for _, name := range names {
		lease, err := b.manager.Begin()
		if err != nil {
			return err
		}
		defer lease.Finish(nil)
		b.success = append(b.success, name)
		if err := b.runner(name); err != nil {
			return err
		}
	}
	return nil
}

func (b *Batch) Successes() []string { return append([]string(nil), b.success...) }
