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

type Response struct {
	Version    string
	StatusCode int
	Message    string
	Headers    map[string]string
	Body       string
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

func formResponse(res *Response) []byte {
	if res.Headers == nil {
		res.Headers = make(map[string]string)
	}
	res.Headers["content-length"] = fmt.Sprintf("%d", len(res.Body))

	if _, ok := res.Headers["content-type"]; !ok {
		res.Headers["content-type"] = "text/plain"
	}

	statusMsg := getStatusMsg(res.StatusCode)

	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", res.StatusCode, statusMsg)

	for key, val := range res.Headers {
		response += fmt.Sprintf("%s: %s\r\n", key, val)
	}

	response += fmt.Sprintf("\r\n%s", res.Body)

	return []byte(response)
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	req, err := parseRequest(conn)
	if err != nil {
		fmt.Println("failed to parse request", err)
		return
	}
	log.Printf("[INFO] %s %s", req.Method, req.Endpoint)

	response := &Response{
		StatusCode: 200,
		Body:       "Hello from my HTTP:)",
	}

	resBytes := formResponse(response)

	conn.Write(resBytes)
}

func getStatusMsg(code int) string {
	switch code {
	case 200:
		return "OK"
	case 404:
		return "Not Found"
	case 400:
		return "Bad Request"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown"
	}
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
