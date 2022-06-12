package router

import (
	"golang.org/x/exp/slices"
	"net/http"
	"wonderful-hand-common/rest"
)

type router struct {
	backend      Backend
	singleRoutes []rest.Route
	groupRoutes  []rest.Group
}

func (r *router) init() {
	r.groupRoutes = []rest.Group{
		rest.NewGroup("/user", nil,
			rest.NewRoute("/register", http.MethodPost, r.handleRegister),
			rest.NewRoute("/login", http.MethodPost, r.handleLogin),
		),
	}
}

func (r *router) SingleRoutes() []rest.Route {
	if r.singleRoutes == nil {
		return make([]rest.Route, 0)
	}
	return slices.Clone(r.singleRoutes)
}

func (r *router) GroupRoutes() []rest.Group {
	if r.groupRoutes == nil {
		return make([]rest.Group, 0)
	}
	return slices.Clone(r.groupRoutes)
}

func New(backend Backend) rest.Router {
	r := &router{backend: backend}
	r.init()
	return r
}
