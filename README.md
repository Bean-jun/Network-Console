# Network Console
A lightweight network debugging tool built with **Vue3 + TypeScript + Vite** and **Golang**. It integrates a visual HTTP request debugger, mock server and HTTP proxy server for convenient API development and network request testing.

## Project Structure
```
network-console/
├── .gitignore                      # Git ignore rules
├── go.mod / go.sum                 # Golang module dependencies
├── index.html                      # Frontend entry HTML
├── main.go                         # Main proxy server entry
├── Makefile                        # Build script
├── package.json / package-lock.json # NPM dependencies
├── tsconfig.json / tsconfig.node.json # TypeScript configuration
├── vite.config.ts                  # Vite build configuration
├── mock-server/                    # Independent Mock server module
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── scripts/                        # Windows startup batch scripts
│   ├── start-mock-server.bat
│   └── start-proxy-server.bat
└── src/                            # Frontend source code
    ├── App.vue                     # Root component
    ├── main.ts                     # Frontend entry file
    ├── vite-env.d.ts               # Vite type declarations
    ├── components/                 # Vue components
    │   ├── RequestBuilder.vue      # HTTP request builder
    │   └── ResponseViewer.vue      # HTTP response viewer
    ├── styles/                     # Global styles
    │   └── main.css
    └── types/                      # Global TypeScript type definitions
        └── index.ts
```

## Features
- ✅ **Visual HTTP Request Builder**: Support all common HTTP methods, custom query params, request headers and request body.
- ✅ **Formatted Response Viewer**: Display HTTP status code, headers, response data, request time and response size.
- ✅ **Built-in HTTP Proxy Server**: Forward network requests to solve cross-domain and network debugging problems.
- ✅ **Modern Tech Stack**: Vue3 + TypeScript + Vite for frontend, Golang + Gin for backend service.
- ✅ **Cross-platform**: Supports Windows, Linux and macOS (recompile binaries for different OS).

## Prerequisites
- **Node.js** 16+ & npm (For frontend development & build)
- **Go** 1.18+ (For Golang service compilation & running)
- Git

## Quick Start

### 1. Clone Repository
```bash
git clone https://github.com/Bean-jun/network-console.git
cd network-console
```

### 2. Run Frontend in Development Mode
```bash
# Install dependencies
npm install

# Start Vite dev server
npm run dev
```
Visit: `http://localhost:3000` (configured in `vite.config.ts`)

### 3. Start Mock Server
#### Option 1: Run batch script (Windows)
```cmd
scripts\start-mock-server.bat
```

#### Option 2: Run with Go command
```bash
cd mock-server
go mod tidy
go run main.go
```

### 4. Start Proxy Server
#### Option 1: Run batch script (Windows)
```cmd
scripts\start-proxy-server.bat
```

#### Option 2: Run with Go command
```bash
go mod tidy
go run main.go
```
Default proxy server port: `7256`

## Build for Production

### Build Frontend Static Files
```bash
npm run build
```
Production files will be generated in the `dist` folder.

### Compile Golang Binaries
#### Windows
```bash
# Compile main proxy server
go build -o network-console.exe main.go

# Compile mock server
cd mock-server
go build -o mock-server.exe main.go
```

#### Linux / macOS
```bash
# Main service
go build -o network-console main.go

# Mock service
cd mock-server
go build -o mock-server main.go
```

## Configuration
- **Frontend Port**: Modify `server.port` in `vite.config.ts`
- **Proxy Server Port**: Change the default port or use flag `--port` when starting the Go program
- **Mock Rules**: Edit routing and response logic in `mock-server/main.go`

## Usage Guide
1. Launch frontend page, fill in request method, URL, headers, params and body.
2. Send request and view formatted response data in the viewer panel.
3. Use the built-in proxy to forward requests and debug network issues.
4. Start mock server to simulate third-party interfaces for local development.

## Notes
- The `files-delete` folder contains temporary test scripts, you can delete it safely.
- `.exe` files are precompiled for Windows only. Recompile binaries for other operating systems.
- The Go backend uses embedded FS to serve frontend `dist` files, you can deploy them together directly.

## Contributing
1. Create a new branch from `main` (e.g. `feature/xxx`, `fix/xxx`).
2. Commit your changes with standard commit messages.
3. Submit a Pull Request.
4. Feel free to create Issues for bugs or feature requests.

## License
This project is open source under the **MIT License**.