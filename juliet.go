package juliet

import (
	"net/http"
)

type contextHandler func(*Context) http.Handler
type contextMiddleware func(ctx *Context, next http.Handler) http.Handler
type contextHandlerFunc func(ctx *Context, w http.ResponseWriter, r *http.Request)

// Chain is a wrapper for a contextMiddleware instance.
// It links contextMiddleware together to create a chain
// to be executed in order.
type Chain struct {
	parent     *Chain
	middleware contextMiddleware
}

// NewChain creates a new contextMiddleware chain.
// As a Chain is just a wrapper to a contextMiddleware instance
// it would make no sense to have an empty chain.
func NewChain(cm contextMiddleware) (chain *Chain) {
	chain = new(Chain)
	chain.middleware = cm
	return
}

// Append wraps a contextMiddleware into a new chain and attach
// it to the current chain of middleware.
func (chain *Chain) Append(cm contextMiddleware) (newChain *Chain) {
	newChain = NewChain(cm)
	newChain.parent = chain
	return newChain
}

// AppendAll appends multiple contextMiddleware to the chain at once
func (chain *Chain) AppendAll(cms ...contextMiddleware) (newChain *Chain) {
	newChain = chain
	for _, cm := range cms {
		newChain = newChain.Append(cm)
	}

	return newChain
}

// Adapt adds context to a middleware so it can be added to the chain
func Adapt(fn func(http.Handler) http.Handler) contextMiddleware {
	return func(ctx *Context, h http.Handler) http.Handler {
		return fn(h)
	}
}

// head returns the top/first middleware of the Chain
func (chain *Chain) head() (head *Chain) {
	// Find the head of the chain
	head = chain
	for head.parent != nil {
		head = head.parent
	}
	return
}

// copy duplicate the whole chain of contextMiddleware
func (chain *Chain) copy() (newChain *Chain) {
	newChain = NewChain(chain.middleware)
	if chain.parent != nil {
		newChain.parent = chain.parent.copy()
	}
	return
}

// AppendChain duplicates a chain and links it to the current chain
// So an append to the old chain don't alter the new one
func (chain *Chain) AppendChain(tail *Chain) (newChain *Chain) {
	// Copy the chain to attach
	newChain = tail.copy()

	// Attach the chain to extend to the new tail
	newChain.head().parent = chain

	// Return the new tail
	return
}

// Then add a contextHandlerFunc to the end of the chain
// and returns a http.Handler compliant ContextHandler
func (chain *Chain) Then(fn contextHandlerFunc) (ch *ContextHandler) {
	ch = NewHandler(chain)
	ch.handler = adaptContextHandlerFunc(fn)
	return
}

// ThenHandler add a http.Handler to the end of the chain
// and returns a http.Handler compliant ContextHandler
func (chain *Chain) ThenHandler(handler http.Handler) (ch *ContextHandler) {
	ch = NewHandler(chain)
	ch.handler = adaptHandler(handler)
	return
}

// ThenHandlerFunc add a http.HandlerFunc to the end of the chain
// and returns a http.Handler compliant ContextHandler
func (chain *Chain) ThenHandlerFunc(fn func(http.ResponseWriter, *http.Request)) (ch *ContextHandler) {
	ch = NewHandler(chain)
	ch.handler = adaptHandlerFunc(fn)
	return
}

// ContextHandler holds a chain and a final handler
// It satisfy the http.Handler interface and can be
// served directly by a net/http server
type ContextHandler struct {
	*Chain
	handler contextHandler
}

// New Handler creates a new handler chain
func NewHandler(chain *Chain) (ch *ContextHandler) {
	ch = new(ContextHandler)
	ch.Chain = chain
	return
}

// Execute build the chain of handlers and execute it
// passing the context along
func (ch *ContextHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := NewContext()

	// Build the context handler chain
	handler := ch.handler(ctx)
	chain := ch.Chain
	for chain != nil {
		handler = chain.middleware(ctx, handler)
		chain = chain.parent
	}

	handler.ServeHTTP(resp, req)
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