package juliet
import (
	"github.com/root-gg/utils"
	"fmt"
)

type Context struct {
	parent *Context
	values map[interface{}]interface{}
}

func NewContext() (ctx *Context) {
	ctx = new(Context)
	ctx.values = make(map[interface{}]interface{})
	return
}

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

func (ctx *Context) Set(key string, val interface{}) {
	ctx.values[key] = val
}

func (ctx *Context) Delete(key string) {
	delete(ctx.values, key)
}

func (ctx *Context) Fork() *Context {
	nc := NewContext()
	nc.parent = ctx
	return nc
}

func (ctx *Context) Dump() {
	fmt.Printf("%+v\n",ctx)
}