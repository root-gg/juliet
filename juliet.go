package juliet

import (
	"net/http"
)

type contextHandler func(*Context) http.Handler
type contextMiddleware func(ctx *Context, next http.Handler) http.Handler
type contextHandlerFunc func(ctx *Context, w http.ResponseWriter, r *http.Request)

type Chain struct {
	parent *Chain
	middleware contextMiddleware
}

func NewChain(cm contextMiddleware) (chain *Chain) {
	chain = new(Chain)
	chain.middleware = cm
	return
}

func (chain *Chain) Append(cm contextMiddleware) (newChain *Chain) {
	newChain = NewChain(cm)
	newChain.parent = chain
	return newChain
}

func (chain *Chain) AppendAll(cms... contextMiddleware) (newChain *Chain) {
	newChain = chain
	for _, cm := range cms {
		newChain = newChain.Append(cm)
	}

	return newChain
}

func (chain *Chain) head() (head *Chain){
	// Find the head of the chain
	head = chain
	for head.parent != nil {
		head = head.parent
	}
	return
}

func (chain *Chain) copy() (newChain *Chain){
	newChain = NewChain(chain.middleware)
	if chain.parent != nil {
		newChain.parent = chain.parent.copy()
	}
	return
}

func (chain *Chain) AppendChain(tail *Chain) (newChain *Chain){
	// Copy the chain to attach
	newChain = tail.copy()

	// Attach the chain to extend to the new tail
	newChain.head().parent = chain

	// Return the new tail
	return
}

func (chain *Chain) Then(fn contextHandlerFunc) (ch *ContextHandler) {
	ch = NewHandler(chain)
	ch.handler = adaptContextHandlerFunc(fn)
	return
}

func (chain *Chain) ThenHandler(handler http.Handler) (ch *ContextHandler) {
	ch = NewHandler(chain)
	ch.handler = adaptHandler(handler)
	return
}

func (chain *Chain) ThenHandlerFunc(fn func(http.ResponseWriter, *http.Request)) (ch *ContextHandler) {
	ch = NewHandler(chain)
	ch.handler = adaptHandlerFunc(fn)
	return
}

type ContextHandler struct {
	*Chain
	handler contextHandler
}

func NewHandler(chain *Chain) (ch *ContextHandler){
	ch = new(ContextHandler)
	ch.Chain = chain
	return
}

func (ch *ContextHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := NewContext()

	// Build the context handler chain
	handler := ch.handler(ctx)
	chain := ch.Chain
	for chain != nil {
		handler = chain.middleware(ctx, handler)
		chain = chain.parent
	}

	handler.ServeHTTP(resp,req)
}

// Adapt a function with the signature
// func(Context, http.ResponseWriter, *http.Request) into a contextHandler
func adaptContextHandlerFunc(fn contextHandlerFunc) contextHandler {
	return func(ctx *Context) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fn(ctx, w, r)
		})
	}
}

// Adapt http.Handler into a contextHandler
func adaptHandler(h http.Handler) contextHandler {
	return func(ctx *Context) http.Handler {
		return h
	}
}

// Adapt a function with the signature
// func(http.ResponseWriter, *http.Request) into a contextHandler
func adaptHandlerFunc(fn func(w http.ResponseWriter, r *http.Request)) contextHandler {
	return adaptHandler(http.HandlerFunc(fn))
}