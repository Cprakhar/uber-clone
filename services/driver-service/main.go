package main

func main() {
	grpcServer := NewgRPCServer(":9100")
	grpcServer.run()
}
