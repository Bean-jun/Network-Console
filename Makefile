run: build
	network-console.exe -port 7256

build: build-web
	go mod tidy
	taskkill /f /im network-console.exe || exit 0
	go build -o network-console.exe

build-web:
	npm i
	del /q /s dist
	npm run build