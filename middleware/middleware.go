package middleware

import (
	"fmt"
	"net/http"

	"github.com/urfave/negroni"
)

var (
	GlobalMiddle GlobalMiddleware
	GroupMiddle  GroupMiddleware
)

type GlobalMiddleware struct {
	negroni.Handler
}

type GroupMiddleware struct {
	negroni.Handler
}

type Middleware interface {
	Validate() negroni.Handler
}

type Handler negroni.Handler

func init() {
	GlobalMiddle = GlobalMiddleware{}
	GroupMiddle = GroupMiddleware{}
}

func (g GlobalMiddleware) Validate() negroni.Handler {
	return &GlobalMiddleware{}
}

func (g GroupMiddleware) Validate() negroni.Handler {
	return &GroupMiddleware{}
}

func (n GlobalMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("This is GlobalMiddleware!")
	next(rw, r)
}

func (g GroupMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("This is GroupMiddleware!")
	next(rw, r)
}
