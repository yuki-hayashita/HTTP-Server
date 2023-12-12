package api

import (
	// "bufio"
	// "fmt"
	"log"
	"strings"
)

const (
	GET = 0
	POST = 1
)

// HTTPRequest represents an HTTP request structure
type HTTPRequest struct {
	Method  string
	Path    string
	Version string
}

// ParseHTTPRequest parses an HTTP request string and returns an HTTPRequest struct
func ParseHTTPRequest(requestString string) (HTTPRequest, error) {
	parts := strings.Split(requestString, " ")
	log.Println(parts)
	// if len(parts) != 3 {
	// 	return HTTPRequest{}, fmt.Errorf("invalid HTTP request format: %s", requestString)
	// }
	log.Println("HERE")
	log.Println(parts)
	log.Println(requestString)
	method := parts[0]
	path := parts[1]
	version := parts[2]

	return HTTPRequest{
		Method:  method,
		Path:    path,
		Version: version,
	}, nil
}