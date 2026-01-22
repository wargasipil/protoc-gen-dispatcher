@echo off
setlocal enabledelayedexpansion

echo Building protoc-gen-dispatcher...

REM Detect GOBIN (preferred)
if defined GOBIN (
    set TARGET_DIR=%GOBIN%
) else (
    REM Fallback to GOPATH/bin
    set TARGET_DIR=\go\bin
)

echo Target directory: %TARGET_DIR%

REM Build the plugin
go build -o protoc-gen-dispatcher.exe .\cmd\generate

if %errorlevel% neq 0 (
    echo Build failed!
    exit /b 1
)

echo Build success.

REM Create target directory if missing
if not exist "%TARGET_DIR%" (
    echo Creating directory %TARGET_DIR%
    mkdir "%TARGET_DIR%"
)

echo Moving binary to %TARGET_DIR%...

move /Y protoc-gen-dispatcher.exe "%TARGET_DIR%"

if %errorlevel% neq 0 (
    echo Failed to move binary!
    exit /b 1
)

echo Done! Installed protoc-gen-dispatcher in %TARGET_DIR%

endlocal
