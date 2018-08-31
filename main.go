package main

func main() {
	server := newGinHTTPServer(ginHTTPServerConfig{
		host: "0.0.0.0",
		port: "8080",
	})
	server.run()
}
