
Juliet is a lightweight middleware chaining helper that pass a Context (map) object
from a middleware to the next one.

It is a rewrite of [Stack](https://github.com/alexedwards/stack) by Alex Edwards 
which is inspired by [Alice](https://github.com/justinas/alice) by Justinas Stankevicius

### Write a ContextMiddleware
```
    // Write a ContextMiddleware
    func middleware(ctx *juliet.Context,w next http.Handler) http.Handler {
        return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
            // Play with the context
            ctx.Set("key", "value")
            
            // Pass the request to the next middleware / handler
            next.ServeHTTP(resp, req)
        })
    }

    // To create a new chain
    chain := juliet.NewChain(middleware)

    // To append a middleware at the end of the chain
    chain = chain.Append(middleware2)
    
    // To append a middleware at the beginning of a chain
    chainWithAuth := juliet.NewChain(authMiddleware).AppendChain(chain)
    
    // Classic middleware without context can be added to the chain using the Adapt function
    func middlewareWithoutContext(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // middleware logic
        next.ServeHTTP(w, r)
      })
    }
    
    // Note that as AppendChain makes a copy chainWithAuth will NOT contain middlewareWithoutContext
    chainWithClassicMiddleware := chain.Append(juliet.Adapt(middlewareWithoutContext))
    
    // Write a ContextHandler
    func handler(ctx *juliet.Context, resp http.ResponseWriter, req *http.Request) {
        // play with context
        value, _ := ctx.Get("key")
       
       // write http response
       resp.Write([]byte(fmt.Sprintf("value is %v\n", value)))
    }
    
    // Execute a middleware chain
    http.Handle("/", chain.Then(handler))
    
    // Classic http.Handler without context
    http.Handle("/404", chain.ThenHandler(http.NotFoundHandler))
    
    // Classic http.HandlerFunc without context
    func pingHandler(w http.ResponseWriter, r *http.Request) {
      w.Write([]byte("pong"))
    }
    http.Handle("/ping", chain.ThenHandlerFunc(pingHandler))
```
