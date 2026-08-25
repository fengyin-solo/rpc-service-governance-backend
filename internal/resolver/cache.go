package resolver

import (
	"errors"

	"rpcgate/internal/catalog"
)

var ErrNotFound = errors.New("route not found")

type Loader interface {
	Build(string) (*catalog.Route, error)
}

type Cache struct {
	loader Loader
	routes map[string]*catalog.Route
}

func New(loader Loader) *Cache {
	return &Cache{loader: loader, routes: make(map[string]*catalog.Route)}
}

func (c *Cache) Load(name string) error {
	route, err := c.loader.Build(name)
	if route != nil {
		c.routes[name] = route
	}
	return err
}

func (c *Cache) Resolve(name, input string) (string, error) {
	route, ok := c.routes[name]
	if !ok {
		return "", ErrNotFound
	}
	return route.Handler(input), nil
}

func (c *Cache) Size() int { return len(c.routes) }
