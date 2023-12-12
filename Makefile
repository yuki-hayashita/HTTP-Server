all:
	go build -o server ./cmd/server 
	go build -o client ./cmd/client
	
clean:
	rm -fv server client