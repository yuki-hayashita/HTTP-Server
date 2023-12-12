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
	// Connect to the server
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

	// Make a GET request for a specific file
	// request := "GET /Users/naytewen/Desktop/HTTP-Server/files/HelloWorld.html HTTP/1.1\r\n\r\n"
	// conn.Write([]byte(request))

	// Read the response from the server
	// buf := make([]byte, 1024)
	// n, err := conn.Read(buf)
	// if err != nil {
	// 	fmt.Println("Error reading response:", err)
	// 	return
	// }

	// response := string(buf[:n])
	// fmt.Println("Server response:")
	// fmt.Println(response)

	// Save the response to a file
	// saveToFile(response, "output.html")
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

func handleInput(conn net.Conn) {
	// need to put this in for loop
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		request, err := reader.ReadString('\n')

		if err != nil {
			log.Fatalln("Read errors: ", err)
		}
		lineSplit := (strings.Split(request, " "))
		switch lineSplit[0] {
		case "GET":
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
