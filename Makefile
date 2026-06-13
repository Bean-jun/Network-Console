# 产物输出目录
BIN_DIR := bin
# 程序名
APP_NAME := network-console

# 默认运行
run: build-windows-amd64
	./$(BIN_DIR)/$(APP_NAME).exe -port 7256

# Windows x86_64 编译
build-windows-amd64: build-web delete-server
	go mod tidy
	set CGO_ENABLED=0&&set GOOS=windows&&set GOARCH=amd64&&go build -ldflags "-w -s" -o $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe

# Linux x86_64
build-linux-amd64: build-web
	go mod tidy
	set CGO_ENABLED=0&&set GOOS=linux&&set GOARCH=amd64&&go build -ldflags "-w -s" -o $(BIN_DIR)/$(APP_NAME)-linux-amd64

# Linux arm64 (树莓派、鲲鹏等)
build-linux-arm64: build-web
	go mod tidy
	set CGO_ENABLED=0&&set GOOS=linux&&set GOARCH=arm64&&go build -ldflags "-w -s" -o $(BIN_DIR)/$(APP_NAME)-linux-arm64

# macOS Intel
build-darwin-amd64: build-web
	go mod tidy
	set CGO_ENABLED=0&&set GOOS=darwin&&set GOARCH=amd64&&go build -ldflags "-w -s" -o $(BIN_DIR)/$(APP_NAME)-darwin-amd64

# macOS M1/M2/M3 arm64
build-darwin-arm64: build-web
	go mod tidy
	set CGO_ENABLED=0&&set GOOS=darwin&&set GOARCH=arm64&&go build -ldflags "-w -s" -o $(BIN_DIR)/$(APP_NAME)-darwin-arm64

# 一键编译全部平台
all-build: build-windows-amd64 build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64
	rd /q /s dist || exit 0

# 仅 Windows 打包压缩
release: build-windows-amd64
	upx -9 $(BIN_DIR)/$(APP_NAME).exe
	rd /q /s dist || exit 0

# 批量 UPX 压缩所有平台二进制（可选）
release-all: all-build
	upx -9 $(BIN_DIR)/*

# Windows 停止进程 + 删除旧 exe
delete-server:
	taskkill /f /im $(APP_NAME).exe || exit 0

# 前端打包
build-web:
	rd /q /s dist || exit 0
	npm i
	npm run build

# 兼容旧 build 别名，默认编译 Windows
build: build-windows-amd64

# 清理所有编译产物与前端打包文件
clean:
	rd /q /s $(BIN_DIR) dist