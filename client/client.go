package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/tracer-demo/util"
)

func main() {
	// init tracer
	closer, err := util.InitTracer()
	if err != nil {
		log.Fatalf("failed to new tracer, err: %v", err)
		return
	}
	defer closer.Close()

	tracer := opentracing.GlobalTracer()
	clientSpan := tracer.StartSpan("client")
	defer clientSpan.Finish()

	url := "http://127.0.0.1:9090/"
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	// Set some tags on the clientSpan to annotate that it's the client span. The additional HTTP tags are useful for debugging purposes.
	ext.SpanKindRPCClient.Set(clientSpan)
	ext.HTTPUrl.Set(clientSpan, url)
	ext.HTTPMethod.Set(clientSpan, http.MethodGet)

	// Inject the client span context into the headers
	tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	resp, _ := http.DefaultClient.Do(req)

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	length, _ := buf.ReadFrom(resp.Body)

	fmt.Println("recv length: ", length)
	fmt.Println("recv content: ", string(buf.Bytes()))
}
