package main

func main() {
	server := newAPIServer(apiServerConfig{
		host:             "0.0.0.0",
		port:             "8080",
		connectionString: "postgres://hellofresh:hellofresh@localhost:5432/hellofresh?sslmode=disable",
	})
	server.run()
}
