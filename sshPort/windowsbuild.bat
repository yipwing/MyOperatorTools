setlocal
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o d:\download\batchfirewall .\sshPort\batchfirewall.go
go build -o d:\download\firewall .\sshPort\firewall.go
endlocal