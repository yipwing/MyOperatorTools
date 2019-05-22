setlocal
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags "-s -w" -o d:\download\backup backup.go
go build -o d:\download\sshd sshd.go
endlocal