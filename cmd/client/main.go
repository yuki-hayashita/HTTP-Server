package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"log"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: must run client with right number of args: <listen-port>")
	}
	addr := ("localhost:" + os.Args[1])
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	handleInput(conn)

}

func extractBody(httpResponse string) (string, error) {
	// Find the index of the first occurrence of "\r\n\r\n" which marks the end of the headers
	headersEnd := strings.Index(httpResponse, "\r\n\r\n")
	if headersEnd == -1 {
		return "", fmt.Errorf("malformed HTTP response: headers end not found")
	}

	// Extract the body content after the headers
	body := httpResponse[headersEnd+4:]
	return body, nil
}

// method that writes the content to file
func writeToFile(content, filename string) {
	content, _ = extractBody(content)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	cmd := exec.Command("open", file.Name())
	cmd.Run()

	fmt.Println("Writing to: ", filename)
}

// method that handles input from command line
func handleInput(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		request, err := reader.ReadString('\n')

		if err != nil {
			log.Fatalln("Read errors: ", err)
		}
		lineSplit := (strings.Split(request, " "))
		switch lineSplit[0] {
		case "GET" :
			conn.Write([]byte(request))
			handleResponse(conn)
		case "quit\n":
			conn.Close()
			os.Exit(0)
			return
		default:
			log.Fatalln("ERROR! INVALID REQUEST!")
		}

	}
}

// method that handles response from server
func handleResponse(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	response := string(buf[:n])
	fmt.Println("Server response:")
	fmt.Println(response)
	writeToFile(response, "output.html")

}
