setlocal
set GOOS=linux
go build -ldflags "-s -w" -o backup backup.go
endlocal