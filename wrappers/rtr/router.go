package rtr

import (
	"context"
	"ewails/wrappers/common"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	cancelRequestKey = "Cancel_Request"
)

type Cancel struct {
	key string
}

var cancelControl = &Cancel{
	key: cancelRequestKey,
}

func CancelRequest(r *http.Request) error {
	value := r.Context().Value(cancelControl)

	cancelRequest, ok := value.(context.CancelFunc)

	if !ok {
		return fmt.Errorf("invalid type of %T", cancelRequest)
	}

	cancelRequest()
	return nil
}

type MiddlewareFn = http.HandlerFunc

type Router struct {
	path       string
	Mux        *http.ServeMux
	middleware []MiddlewareFn
}

func (router *Router) createHandler(options HandlerOptions, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rCtxCopy, cancelFn := context.WithCancel(r.Context())
		*r = *r.WithContext(context.WithValue(rCtxCopy, cancelControl, cancelFn))

		for _, fn := range router.middleware {
			fn(w, r)

			if r.Context().Err() != nil {
				return
			}
		}

		for _, fn := range options.middleware {
			fn(w, r)

			if r.Context().Err() != nil {
				return
			}
		}

		handler(w, r)
	}
}

// TODO: make everything an option?
func New(path string) *Router {
	if path[len(path)-1] != '/' {
		path += "/"
	}

	return &Router{
		Mux:  http.NewServeMux(),
		path: path,
	}
}

func (r *Router) Use(path string, router *Router) {
	basePath := common.ComposePath(r.path, path)
	fullPath := common.ComposePath(basePath, router.path)

	prefix := basePath[:len(basePath)-1]

	r.Mux.Handle(fullPath, http.StripPrefix(prefix, router.Mux))
}

func (r *Router) UseHandler(path string, handler http.Handler) {
	basePath := common.ComposePath(r.path, path)

	prefix := basePath[:len(basePath)-1]

	h := r.createHandler(HandlerOptions{}, handler.ServeHTTP)

	r.Mux.Handle(basePath, http.StripPrefix(prefix, h))
}

func (r *Router) UseMux(path string, mux *http.ServeMux) {
	fullPath := common.ComposePath(r.path, path)

	prefix := fullPath[:len(fullPath)-1]

	r.Mux.Handle(fullPath, http.StripPrefix(prefix, mux))
}

func (r *Router) UseMiddleware(middlewareFns ...MiddlewareFn) {
	r.middleware = append(r.middleware, middlewareFns...)
}

func (r *Router) RegisterPath(method string, path string, handler http.HandlerFunc, optFns ...HandlerOptsFn) {
	options := common.CreateOptions(&HandlerOptions{}, optFns...)

	fullHandler := r.createHandler(*options, handler)
	fullPath := method + " " + common.ComposePath(r.path, path)

	r.Mux.HandleFunc(fullPath, fullHandler)
}

func (r *Router) Get(path string, handler http.HandlerFunc, optFns ...HandlerOptsFn) {
	r.RegisterPath("GET", path, handler, optFns...)
}

func (r *Router) Post(path string, handler http.HandlerFunc, optFns ...HandlerOptsFn) {
	r.RegisterPath("POST", path, handler, optFns...)
}

func (r *Router) Put(path string, handler http.HandlerFunc, optFns ...HandlerOptsFn) {
	r.RegisterPath("PUT", path, handler, optFns...)
}

func (r *Router) Patch(path string, handler http.HandlerFunc, optFns ...HandlerOptsFn) {
	r.RegisterPath("PATCH", path, handler, optFns...)
}

func (r *Router) Delete(path string, handler http.HandlerFunc, optFns ...HandlerOptsFn) {
	r.RegisterPath("DELETE", path, handler, optFns...)
}

func (r *Router) ToServer(addr string) *http.Server {
	return &http.Server{
		Addr:         addr,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      r.Mux,
	}

}

// TODO: Use optFuncs to provide options and allow custom server override
func (r *Router) Run(addr string) {
	server := r.ToServer(addr)

	fmt.Println("starting server on " + addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
