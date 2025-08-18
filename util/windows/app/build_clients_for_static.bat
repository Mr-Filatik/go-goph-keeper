cd ..\..\..\

echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o static/client-linux-amd64 cmd/client/main.go

echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o static/client-macos-amd64 cmd/client/main.go

echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o static/client-windows-amd64.exe cmd/client/main.go

echo "Done! Binaries in /static"

pause