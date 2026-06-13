taskkill /f /im mock-server.exe
cd mock-server
go mod tidy
go build -o ..\mock-server.exe main.go
cd ..
mock-server.exe
