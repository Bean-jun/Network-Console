run: build
	network-console.exe -port 7256

delete-dist:
	del /q /s dist

delete-server:
	taskkill /f /im network-console.exe || exit 0
	del /q network-console.exe

build: build-web delete-server
	go mod tidy
	go build -ldflags "-w -s" -o network-console.exe

build-web: delete-dist
	npm i
	npm run build

release: build
	upx -9 network-console.exe