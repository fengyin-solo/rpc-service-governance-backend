package noderegistry

import "sync"

type Registry struct {
	mu    sync.RWMutex
	nodes map[string]int
}

func New() *Registry { return &Registry{nodes: make(map[string]int)} }

func (r *Registry) Put(address string, weight int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[address] = weight
}

// Delete 下线指定地址的节点。与 Put 一样只写注册表内部状态，
// 不会影响此前已通过 Snapshot 返回的快照副本。
func (r *Registry) Delete(address string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.nodes, address)
}

// Snapshot 返回当前节点权重的一个拷贝。
// 必须返回独立副本而非内部 map 的引用，否则正在刷新负载均衡的一轮
// 在“已拿到节点列表”与“计算总权重”之间若发生并发注册/下线，
// 会被直接改写这一轮本应使用的开始时快照，导致总权重串入新节点。
func (r *Registry) Snapshot() map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]int, len(r.nodes))
	for k, v := range r.nodes {
		out[k] = v
	}
	return out
}
