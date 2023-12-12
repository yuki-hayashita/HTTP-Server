package main

import (
	"fmt"
	api "http-server/pkg"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var (
	HTMLDirectory = "/"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: must run server with right number of args: ./server <listen-port> <path-to-directory>")
	}
	portNumber := os.Args[1]
	HTMLDirectory = strings.TrimSuffix(os.Args[2], "/")
	addr := ":" + portNumber
	listener, err := net.Listen("tcp", addr)

	_, err = directoryExists(HTMLDirectory)
	if err != nil {
		log.Println("Error checking directory existence:", err)
		return
	}

	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port:", portNumber)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Connection Established with a client!")
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading request:", err)
		return
	}


	request, _ := api.ParseHTTPRequest(string(buf))

	if request.Method != "GET" {
		log.Println(request.Method)
		log.Println("INVALID REQUEST!")
		return
	}

	serveContent(conn, request.Path)
}

func serveContent(conn net.Conn, path string) {
	filePath := HTMLDirectory + path
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("COULD NOT FIND FILE!")
		return
	}
	defer file.Close()

	// Send HTTP header

	// Copy the file content to the response writer
	// io.Copy(conn, file)
	serveHTMLFile(conn, file)
}

func serveHTMLFile(conn net.Conn, file *os.File) error {
	// Read the content of the file
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	content := make([]byte, fileInfo.Size())
	_, err = file.Read(content)
	if err != nil && err != io.EOF {
		return err
	}

	// Create the HTTP response
	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/html\r\n"+
			"Content-Length: %d\r\n"+
			"\r\n"+
			"%s",
		fileInfo.Size(),
		content,
	)

	// Write the HTTP response to the connection
	_, err = conn.Write([]byte(response))
	if err != nil {
		return err
	}

	return nil
}

func directoryExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // Directory does not exist
		}
		return false, err // Other error (e.g., permission issues)
	}

	// The directory exists
	return true, nil
}
