run: build
	network-console.exe -port 7256

release: build
	upx -9 network-console.exe

delete-server:
	taskkill /f /im network-console.exe || exit 0
	del /q network-console.exe

build-web:
	del /q /s dist
	npm i
	npm run build

build: build-web delete-server
	go mod tidy
	go build -ldflags "-w -s" -o network-console.exe
	del /q /s dist
