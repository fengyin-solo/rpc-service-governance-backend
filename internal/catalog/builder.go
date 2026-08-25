package catalog

import "fmt"

type Route struct {
	Name    string
	Handler func(string) string
}

type Builder struct{ beforeHandler func() }

func NewBuilder(beforeHandler func()) *Builder { return &Builder{beforeHandler: beforeHandler} }

func (b *Builder) Build(name string) (route *Route, err error) {
	route = &Route{Name: name}
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("build route: %v", recovered)
		}
	}()
	if b.beforeHandler != nil {
		b.beforeHandler()
	}
	route.Handler = func(input string) string { return name + ":" + input }
	return route, nil
}
