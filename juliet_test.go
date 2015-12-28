package juliet

import (
	"net/http"
	"fmt"
	"testing"
//	"github.com/gorilla/mux"
	"net/http/httptest"
	"log"
	"io/ioutil"
	"github.com/gorilla/mux"
)

func appMiddleware1(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request){
		fmt.Println("app middlerware1")
		ctx.Set("value",1)
		next.ServeHTTP(resp,req)
	})
}

func appMiddleware2(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request){
		fmt.Printf("vars %v",mux.Vars(req))
		fmt.Println("app middlerware2")
		ctx.Set("value",2)
		next.ServeHTTP(resp,req)
	})
}

func appMiddleware3(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request){
		fmt.Println("app middlerware3")
		ctx.Set("value",3)
		next.ServeHTTP(resp,req)
	})
}

func appHandler(ctx *Context, resp http.ResponseWriter, req *http.Request){
	fmt.Println("app handler")
	value, _ := ctx.Get("value")
	resp.Write([]byte(fmt.Sprintf("values is %v\n",value)))
}

func appHandlerFunc(resp http.ResponseWriter, req *http.Request){
	fmt.Println("app handler func")
	resp.Write([]byte("yo"))
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

func Test(t *testing.T){
	stack := NewChain(appMiddleware1)
	stack2 := stack.Append(appMiddleware2)
//	stack3 := stack2.Append(appMiddleware3)

//	stack4 := NewChain(appMiddleware2).AppendChain(stack3)
//	stack5 := stack3.AppendChain(stack3)
//
	r := mux.NewRouter()
//	fmt.Println(serveAndRequest(stack.Then(appHandler)))
	r.Handle("/2/{id}", stack2.Then(appHandler))
//	r.Handle("/3", stack3.Then(appHandler))
//	r.Handle("/4", stack4.Then(appHandler))
//	r.Handle("/5", stack5.Then(appHandler))
//	r.Handle("/6", stack5.ThenHandlerFunc(appHandlerFunc))
//	r.Handle("/7", stack5.ThenHandler(http.NotFoundHandler()))
//	r.Handle("/8/{id}", stack5.ThenHandler(http.NotFoundHandler()))
////	http.Handle("/", appMiddleware(http.HandlerFunc(appHandler)))
	http.ListenAndServe(":3000",r)
}