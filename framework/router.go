package framework

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	Instance chi.Router
}

var singleton *Router
var routerOnce sync.Once

func GetRouter() *Router {
	if singleton == nil {
		routerOnce.Do(func() {
			singleton = &Router{Instance: chi.NewRouter()}
		})
	}
	return singleton
}

func (r *Router) AddPostHandler(path string, handlerFunc http.HandlerFunc) {
	r.Instance.Post(path, handlerFunc)
}

func (r *Router) AddGetHandler(path string, handlerFunc http.HandlerFunc) {
	r.Instance.Get(path, handlerFunc)
}
