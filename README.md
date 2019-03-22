# Quick Start

1. Run jaeger locally via docker:

	`docker run -d -p 5775:5775/udp -p 16686:16686 jaegertracing/all-in-one:latest`

2. Run the tracer demo:

	`go run server/server.go`
	
	`go run client/client.go`

3. View your trace at localhost:16686