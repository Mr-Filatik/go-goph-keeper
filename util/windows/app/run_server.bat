@echo off

setlocal enabledelayedexpansion

cd ..\..\..\

for /f "tokens=1,* delims==" %%a in (util\env\local.env) do (
    set %%a=%%b
)

set "ARGS="

if defined CRYPTO_JWT_KEY (
    set "ARGS=!ARGS! -crypto-jwt-key=!CRYPTO_JWT_KEY!"
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

go run cmd\server\main.go !ARGS!

endlocal

pause