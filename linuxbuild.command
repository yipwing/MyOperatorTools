# make sure your system has go environment and has setting the gopath.

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./backup ./backup.go
