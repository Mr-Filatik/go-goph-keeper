@echo off

setlocal enabledelayedexpansion

cd ..\..\..\

rem --- READ LDFLAGS IN PACKAGE FROM GIT AND SYSTEM ---

set "PACKAGE_PATH=github.com/mr-filatik/go-goph-keeper/internal/server"

echo PACKAGE_PATH: !PACKAGE_PATH!

set "LDFLAGS="

if not defined BUILD_DATE (
    for /f %%i in ('powershell -NoProfile -Command "Get-Date -Format yyyy-MM-dd"') do set BUILD_DATE=%%i
)

if not defined BUILD_COMMIT (
    for /f %%i in ('git rev-parse --short HEAD 2^>NUL') do set BUILD_COMMIT=%%i
)

if not defined BUILD_VERSION (
    for /f %%i in ('git describe --tags --always --dirty 2^>NUL') do set BUILD_VERSION=%%i
)

if defined BUILD_VERSION (
    set "LDFLAGS=!LDFLAGS!-X !PACKAGE_PATH!.buildVersion=!BUILD_VERSION!"
)

if defined BUILD_DATE (
    set "LDFLAGS=!LDFLAGS! -X !PACKAGE_PATH!.buildDate=!BUILD_DATE!"
)

if defined BUILD_COMMIT (
    set "LDFLAGS=!LDFLAGS! -X !PACKAGE_PATH!.buildCommit=!BUILD_COMMIT!"
)

if defined LDFLAGS (
    set "LDFLAGS=-ldflags ^"!LDFLAGS!^""
)

echo LDFLAGS: !LDFLAGS!

rem --- READ ARGS FROM .ENV FILE ---

set "ENV_FILE_PATH=util\env\local.env"

echo ENV_FILE_PATH: !ENV_FILE_PATH!

set "ARGS="

if exist "!ENV_FILE_PATH!" (
    for /f "tokens=1,* delims==" %%a in (!ENV_FILE_PATH!) do (
        set %%a=%%b
    )
)

if defined CRYPTO_JWT_KEY (
    set "ARGS=!ARGS!-crypto-jwt-key=!CRYPTO_JWT_KEY!"
)

if defined DATABASE_CONNECTION_STRING (
    set "ARGS=!ARGS! -database=!DATABASE_CONNECTION_STRING!"
)

if defined HASH_KEY (
    set "ARGS=!ARGS! -hash-key=!HASH_KEY!"
)

if defined SERVER_ADDRESS (
    set "ARGS=!ARGS! -address=!SERVER_ADDRESS!"
)

echo ARGS: !ARGS!

rem --- RUN APP WITH LDFLAGS AND ARGS ---

set "APP_PATH=cmd\server\main.go"

echo APP_PATH: !APP_PATH!
echo.

go run !LDFLAGS! !APP_PATH! !ARGS!

endlocal

@echo on

pause