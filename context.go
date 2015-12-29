package juliet

import (
	"fmt"
)

// Context is a recursive data structure to hold values
// A Fork can access ancestors values, it might override
// them locally. But it can not alter ancestors state.
type Context struct {
	parent *Context
	values map[interface{}]interface{}
}

// NewContext creates a new context instance
func NewContext() (ctx *Context) {
	ctx = new(Context)
	ctx.values = make(map[interface{}]interface{})
	return
}

// Get returns the value matching the key from the context.
func (ctx *Context) Get(key string) (interface{}, bool) {
	if value, ok := ctx.values[key]; ok {
		return value, true
	} else {
		if ctx.parent != nil {
			return ctx.parent.Get(key)
		}
	}
	return nil, false
}

// Set adds a value to the context or overrides a parent value
func (ctx *Context) Set(key string, val interface{}) {
	ctx.values[key] = val
}

// Delete removes a key from the context
// This has no effect if the key is defined in
// a parent context.
func (ctx *Context) Delete(key string) {
	delete(ctx.values, key)
}

// Fork creates a new child context
func (ctx *Context) Fork() *Context {
	nc := NewContext()
	nc.parent = ctx
	return nc
}

// String returns a context string representation
func (ctx *Context) String() (str string) {
	if ctx.parent != nil {
		str += ctx.parent.String()
	}
	for key, value := range ctx.values {
		str += fmt.Sprintf("%v => %v\n", key, value)
	}
	return
}
