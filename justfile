set shell := ["fish", "-c"]

build:
    go build -o dist/main cmd/main.go

mcp-debugger:
    bunx @modelcontextprotocol/inspector ./dist/main
    
test:
    go test ./internal/...