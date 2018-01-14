@echo off
SETLOCAL
set VERSION=0.0.2

echo "Compiling windows 64bit"
set GOOS=windows
set GOARCH=amd64
REM go build -ldflags "-H windowsgui -X main.Version=%VERSION%" -o %VERSION%/mstt-client-%GOOS%-%GOARCH%-%VERSION%.exe

echo "Compiling Linux 64bit"
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-X main.Version=%VERSION%" -o %VERSION%/mstt-client-%GOOS%-%GOARCH%-%VERSION%

echo "Compiling Freebsd 64bit"
set GOOS=freebsd
set GOARCH=amd64
REM go build -ldflags "-X main.Version=%VERSION%" -o %VERSION%/mstt-client-%GOOS%-%GOARCH%-%VERSION%

echo "Compiling OSX 64bit"
set GOOS=darwin
set GOARCH=amd64
REM go build -ldflags "-X main.Version=%VERSION%" -o %VERSION%/mstt-client-%GOOS%-%GOARCH%-%VERSION%
REM curl -F "file=@mstt-client-windows-%VERSION%.exe" http://192.168.20.149/upload.php