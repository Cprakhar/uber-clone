package main

import "github.com/cprakhar/uber-clone/shared/env"

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8080")
)

func main() {
	httpServer := NewHTTPServer(httpAddr)
	go httpServer.run()
}
