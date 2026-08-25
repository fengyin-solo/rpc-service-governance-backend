package routecache

type Cache struct{ entries map[string][]byte }

func New() *Cache { return &Cache{entries: make(map[string][]byte)} }

func (c *Cache) Put(key string, route []byte) {
	c.entries[key] = route
}

func (c *Cache) Get(key string) []byte {
	return c.entries[key]
}
