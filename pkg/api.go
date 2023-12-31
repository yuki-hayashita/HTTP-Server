package api

import (
	// "bufio"
	// "fmt"
	"fmt"
	"strings"
)

const (
	GET = 0
	POST = 1
)

// HTTPRequest struct
type HTTPRequest struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    string
}

// method to convert request from client into HTTPRequest struct
func ParseHTTPRequest(requestString string) (HTTPRequest, error) {
	request := HTTPRequest{
		Headers: make(map[string]string),
	}

	lines := strings.Split(requestString, "\r\n")

	// parse the requests
	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) != 3 {
		return HTTPRequest{}, fmt.Errorf("invalid HTTP request format: %s", requestString)
	}
	request.Method = requestLine[0]
	request.Path = requestLine[1]
	request.Version = requestLine[2]

	// parse the headers
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			// Empty line marks the end of headers
			break
		}
		headerParts := strings.SplitN(lines[i], ": ", 2)
		if len(headerParts) == 2 {
			request.Headers[headerParts[0]] = headerParts[1]
		}
	}

	if len(lines) > 0 {
		request.Body = lines[len(lines)-1]
	}

	return request, nil
}