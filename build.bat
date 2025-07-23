@echo off
REM Build script for Go applications
REM Creates optimized binaries for multiple platforms
REM Usage: build.bat [base-filename]
REM Example: build.bat my-app

REM Set base filename from argument or use default
if "%1"=="" (
    set BASE_NAME=bubblegum-physics-sim
) else (
    set BASE_NAME=%1
)

echo Building %BASE_NAME% for multiple platforms...

REM Create build directory
mkdir builds 2>nul

REM Build for Linux (64-bit)
echo Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o builds/%BASE_NAME%-linux .

REM Build for Windows (64-bit)
echo Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o builds/%BASE_NAME%.exe .

REM Build for macOS Intel (64-bit)
echo Building for macOS Intel (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -o builds/%BASE_NAME%-macos-intel .

REM Build for macOS Apple Silicon (arm64)
echo Building for macOS Apple Silicon (arm64)...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-s -w" -o builds/%BASE_NAME%-macos-arm .

REM Build for current platform (local testing)
echo Building for current platform...
set GOOS=
set GOARCH=
go build -ldflags="-s -w" -o builds/%BASE_NAME%-local.exe .

echo Build complete! Binaries are in the 'builds' directory:
dir builds\ 