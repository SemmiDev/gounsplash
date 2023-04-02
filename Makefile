windows64:
	GOOS=windows GOARCH=amd64 go build -o main.exe main.go
windows32:
	GOOS=windows GOARCH=386 go build -o main.exe main.go
linux:
	GOOS=linux GOARCH=amd64 go build -o main main.go
