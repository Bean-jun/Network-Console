taskkill /f /im network-console.exe
go mod tidy
go build -o network-console.exe main.go
start network-console.exe -port 7256
