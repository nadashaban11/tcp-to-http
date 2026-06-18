package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Request struct {
	Method   string
	Endpoint string
	Version  string
	Headers  map[string]string
}

func parseRequest(conn net.Conn) (*Request, error) {
	reader := bufio.NewReader(conn)

	firstLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	firstLine = strings.TrimSpace(firstLine)
	reqParts := strings.Split(firstLine, " ")
	if len(reqParts) != 3 {
		return nil, fmt.Errorf("not allowed request line")
	}

	request := &Request{
		Method:   reqParts[0],
		Endpoint: reqParts[1],
		Version:  reqParts[2],
		Headers:  make(map[string]string),
	}

	// to get headers

	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		headerLine = strings.TrimSpace(headerLine)

		if headerLine == "" {
			break
		}

		headerParts := strings.SplitN(headerLine, ":", 2)
		key := strings.TrimSpace(headerParts[0])
		val := strings.TrimSpace(headerParts[1])

		request.Headers[key] = val
	}
	return request, nil
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	req, err := parseRequest(conn)
	if err != nil {
		fmt.Println("failed to parse request", err)
		return
	}
	log.Printf("[INFO] %s %s", req.Method, req.Endpoint)

	body := "hello from my HTTP:)"
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
	conn.Write([]byte(response))
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("refused to connect to port 8080", err)
	}
	defer listener.Close()

	log.Printf("tcp is listening on http://localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("failed to accept connection ", err)
			continue
		}
		go handleConn(conn)
	}

}
