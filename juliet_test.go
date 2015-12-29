package juliet

import (
	"fmt"
	"net/http"
	"testing"
	"io/ioutil"
	"log"
	"net/http/httptest"
)

func appMiddleware1(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		fmt.Println("app middlerware1")
		ctx.Set("1", true)
		ctx.Set("last", 1)
		next.ServeHTTP(resp, req)
	})
}

func appMiddleware2(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		fmt.Println("app middlerware2")
		ctx.Set("2", true)
		ctx.Set("last", 2)
		next.ServeHTTP(resp, req)
	})
}

func appMiddleware3(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		fmt.Println("app middlerware3")
		ctx.Set("3", true)
		ctx.Set("last", 3)
		next.ServeHTTP(resp, req)
	})
}

func appHandler(ctx *Context, resp http.ResponseWriter, req *http.Request) {
	fmt.Println("app handler")
	value, _ := ctx.Get("last")
	resp.Write([]byte(fmt.Sprintf("value is %v\n", value)))
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

func Test(t *testing.T) {
	chain := NewChain(appMiddleware1)
	chain2 := chain.Append(appMiddleware2)
	http.ListenAndServe(":3000", chain2.Then(appHandler))
}
