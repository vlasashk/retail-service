package muxer

import "net/http"

type Middleware func(http.Handler) http.Handler

type MyMux struct {
	*http.ServeMux
	middlewares []Middleware // глобальные
}

func NewMyMux() *MyMux {
	return &MyMux{
		ServeMux:    http.NewServeMux(),
		middlewares: make([]Middleware, 0),
	}
}

// Use создает список глобальных middleware в порядке, когда самый первый оборачивает все хендлеры, а самый последний оборачивает дефолтный
func (mux *MyMux) Use(m ...Middleware) {
	if mux.middlewares == nil {
		mux.middlewares = make([]Middleware, 0, len(m))
	}
	mux.middlewares = append(mux.middlewares, m...)
}

func (mux *MyMux) Chain() http.Handler {
	return chain(mux, mux.middlewares...)
}

// chain оборачивает middleware в обратном порядке, так что первый в списке оборачивает все последующие
func chain(h http.Handler, m ...Middleware) http.Handler {
	if len(m) < 1 {
		return h
	}
	newHandler := h
	for i := len(m) - 1; i >= 0; i-- {
		newHandler = m[i](newHandler)
	}
	return newHandler
}
