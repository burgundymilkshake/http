package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	// 1. Bind to a TCP port to listen for raw network streams
	port := ":8080"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to bind to port %s: %v", port, err)
	}
	defer listener.Close()

	fmt.Printf("Custom HTTP Server running from scratch on http://localhost%s\n", port)

	for {
		// 2. Accept incoming TCP connections from clients (like browsers or curl)
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// 3. Handle each connection concurrently in a goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create a buffered reader to read the network stream line-by-line
	reader := bufio.NewReader(conn)

	// 4. Parse the HTTP Request Line (e.g., "GET /hello HTTP/1.1")
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading request line: %v", err)
		return
	}

	// Clean up carriage returns (\r\n)
	requestLine = strings.TrimSpace(requestLine)
	if requestLine == "" {
		return
	}

	// Split the request line into parts
	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 {
		log.Printf("Malformed HTTP request line: %s", requestLine)
		return
	}

	method := parts[0]
	path := parts[1]
	protocol := parts[2]

	fmt.Printf("[Request] %s %s (%s)\n", method, path, protocol)

	// 5. Consume HTTP Headers (read until an empty line "\r\n" is encountered)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.TrimSpace(line) == "" {
			break // End of HTTP headers reached
		}
	}

	// 6. Handle Routing and Formulate the HTTP/1.1 Response
	var responseBody string
	var statusCode string

	switch path {
	case "/":
		statusCode = "200 OK"
		responseBody = "<h1>Welcome to my From-Scratch HTTP Server!</h1>"
	case "/hello":
		statusCode = "200 OK"
		responseBody = fmt.Sprintf("<h1>Hello! You requested the path: %s</h1>", path)
	default:
		statusCode = "404 Not Found"
		responseBody = "<h1>404 Page Not Found</h1>"
	}

	// 7. Write the raw HTTP/1.1 response structure back to the TCP socket
	// HTTP protocol requires exact sequence: Status-Line -> Headers -> Empty Line -> Body
	response := fmt.Sprintf(
		"HTTP/1.1 %s\r\n"+
			"Content-Type: text/html; charset=utf-8\r\n"+
			"Content-Length: %d\r\n"+
			"Connection: close\r\n"+
			"\r\n"+ // Crucial blank line dividing headers from body
			"%s",
		statusCode, len(responseBody), responseBody,
	)

	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
