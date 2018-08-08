all: testclient testserver

testclient: client/main.go
	go build -o testclient ./client/...

testserver: server/main.go
	go build -o testserver ./server/...
