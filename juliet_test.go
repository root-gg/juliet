package juliet

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func middleware1(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "middlerware1->")
		next.ServeHTTP(resp, req)
	})
}

func middleware2(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "middlerware2->")
		next.ServeHTTP(resp, req)
	})
}

func middleware3(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "middlerware3->")
		next.ServeHTTP(resp, req)
	})
}

func thirdPartyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "thirdPartyMiddleware->")
		next.ServeHTTP(resp, req)
	})
}

func handler(ctx *Context, resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "handler")
}

func handlerFunc(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "handler")
}

func serveAndRequest(h http.Handler) string {
	ts := httptest.NewServer(h)
	defer ts.Close()
	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(resBody)
}

func TestEmptyChain(t *testing.T) {
	str := serveAndRequest(NewChain().Then(handler))
	expected := "handler"
	if str != expected {
		t.Fatalf("Invalid output %s, expected %s", str, expected)
	}
}

func TestConstructor(t *testing.T) {
	str := serveAndRequest(NewChain(middleware1, middleware2, middleware3).Then(handler))
	expected := "middlerware1->middlerware2->middlerware3->handler"
	if str != expected {
		t.Fatalf("Invalid output %s, expected %s", str, expected)
	}
}

func TestAppend(t *testing.T) {
	str := serveAndRequest(NewChain().Append(middleware1, middleware2, middleware3).Then(handler))
	expected := "middlerware1->middlerware2->middlerware3->handler"
	if str != expected {
		t.Fatalf("Invalid output %s, expected %s", str, expected)
	}
}

func TestAdapt(t *testing.T) {
	str := serveAndRequest(NewChain(Adapt(thirdPartyMiddleware)).Then(handler))
	expected := "thirdPartyMiddleware->handler"
	if str != expected {
		t.Fatalf("Invalid output %s, expected %s", str, expected)
	}
}

func TestAppendChain(t *testing.T) {
	str := serveAndRequest(NewChain(middleware1).AppendChain(NewChain(middleware2, middleware3)).Then(handler))
	expected := "middlerware1->middlerware2->middlerware3->handler"
	if str != expected {
		t.Fatalf("Invalid output %s, expected %s", str, expected)
	}
}

func TestThenHandler(t *testing.T) {
	str := serveAndRequest(NewChain().ThenHandler(http.HandlerFunc(handlerFunc)))
	expected := "handler"
	if str != expected {
		t.Fatalf("Invalid output %s, expected %s", str, expected)
	}
}

func TestThenHandlerFunc(t *testing.T) {
	str := serveAndRequest(NewChain().ThenHandlerFunc(handlerFunc))
	expected := "handler"
	if str != expected {
		t.Fatalf("Invalid output %s, expected %s", str, expected)
	}
}

func TestContext(t *testing.T) {

	contextMiddleware1 := func(ctx *Context, next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			ctx.Set("value", "contextMiddleware1->")
			next.ServeHTTP(resp, req)
		})
	}

	contextMiddleware2 := func(ctx *Context, next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if value, ok := ctx.Get("value"); ok {
				ctx.Set("value", value.(string)+"contextMiddleware2->")
			}
			next.ServeHTTP(resp, req)
		})
	}

	contextHandler := func(ctx *Context, resp http.ResponseWriter, req *http.Request) {
		if value, ok := ctx.Get("value"); ok {
			fmt.Fprintf(resp, value.(string)+"contextHandler")
		}

	}

	str := serveAndRequest(NewChain(contextMiddleware1, contextMiddleware2).Then(contextHandler))
	expected := "contextMiddleware1->contextMiddleware2->contextHandler"
	if str != expected {
		t.Fatalf("Invalid output %s, expected %s", str, expected)
	}
}
