package main

func main() {
	gRPCServer := NewgRPCServer(":9000")
	gRPCServer.run()
}
