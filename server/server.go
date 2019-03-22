package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/tracer-demo/util"
)

func sayhello(w http.ResponseWriter, r *http.Request) {
	log.Println("recv request from client: ", r.Host)

	tracer := opentracing.GlobalTracer()
	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	serverSpan := tracer.StartSpan("server", ext.RPCServerOption(spanCtx))
	defer serverSpan.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), serverSpan)

	work(ctx)

	fmt.Fprintf(w, "Hello!") //这个写入到w的是输出到客户端的
}

func work(ctx context.Context) {
	tracer := opentracing.GlobalTracer()
	// parentSpan存储到context中, 这里取出来可以直接拿到, 这样就不需要传递span参数了
	// 这种用法只适合childOf方式
	workASpan, _ := opentracing.StartSpanFromContext(ctx, "workA")
	// workASpan := tracer.StartSpan("workA", opentracing.ChildOf(parentSpan.Context()))
	defer workASpan.Finish()

	fmt.Println("this is work method stage A")
	time.Sleep(time.Second * 1)

	workBSpan := tracer.StartSpan("workB", opentracing.FollowsFrom(workASpan.Context()))
	defer workBSpan.Finish()

	fmt.Println("this is work method stage B")
	time.Sleep(time.Second * 1)
}

func main() {
	// init tracer
	closer, err := util.InitTracer()
	if err != nil {
		log.Fatalf("failed to new tracer, err: %v", err)
		return
	}
	defer closer.Close()

	http.HandleFunc("/", sayhello) //设置访问的路由
	log.Println("start server...")
	if err = http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
