setlocal
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o d:\download\batchssh batchssh.go
endlocal