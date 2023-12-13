package main

import (
	"fmt"
	api "http-server/pkg"
	"io"
	"strconv"

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
	if err != nil{
		log.Println("Error creating listener!", err)
		return
	}

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

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading request:", err)
		return
	}

	request, _ := api.ParseHTTPRequest(string(buf))

	if request.Method == "GET" {
		serveContent(conn, request.Path)
	} else if request.Method == "POST" {
		log.Println(request.Headers["Content-Length"])

		if request.Body[0:5] != "name=" {
			log.Println("Invalid val!")
		}
		length, err := strconv.Atoi(request.Headers["Content-Length"])
		if err != nil {
			log.Println("ERROR IN CONVERSION")
			return
		}
		addListItemToFile(request.Path, request.Body[5:length])
		postResp(conn)
	} else {
		log.Println(request.Method)
		log.Println("INVALID REQUEST!")
	}

}

// method that returns the content of an HTML page
func serveContent(conn net.Conn, path string) {
	filePath := HTMLDirectory + path
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("COULD NOT FIND FILE!")
		return
	}
	defer file.Close()

	getResp(conn, file)
}

// method that returns a response for GET request
func getResp(conn net.Conn, file *os.File) error {
	defer conn.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	content := make([]byte, fileInfo.Size())
	_, err = file.Read(content)
	if err != nil && err != io.EOF {
		return err
	}


	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/html\r\n"+
			"Content-Length: %d\r\n"+
			"\r\n"+
			"%s",
		fileInfo.Size(),
		content,
	)

	_, err = conn.Write([]byte(response))
	if err != nil {
		return err
	}

	return nil
}

// method that returns a response for POSE request
func postResp(conn net.Conn) error {
	defer conn.Close()

	content := "{'status' : 'success'}"

	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/html\r\n"+
			"Content-Length: %d\r\n"+
			"\r\n"+
			"%s",
		len(content),
		content,
	)

	_, err := conn.Write([]byte(response))
	if err != nil {
		return err
	}

	return nil
}

// method that checks if directory exists
func directoryExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil 
		}
		return false, err 
	}

	return true, nil
}

// method that adds name to list of signup 
func addListItemToFile(filePath, name string) {
	filePath = HTMLDirectory + filePath
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("COULD NOT FIND FILE!")
		return 
	}
	defer file.Close()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return 
	}
	insertIndex := strings.Index(string(content), "</ol>")
	if insertIndex == -1 {
		log.Println("couldn't find </ol> tag in HTML file")
		return 
	}

	newListItem := fmt.Sprintf("        <li>%s</li>\n", name)

	// inserts name into list on sign-up page
	modifiedContent := string(content[:insertIndex]) + newListItem + string(content[insertIndex:])

	err = os.WriteFile(filePath, []byte(modifiedContent), 0644)
	if err != nil {
		return
	}

	 
}
