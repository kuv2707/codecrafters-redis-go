@echo off
REM Create a temporary file
set tmpFile=%temp%\%random%.tmp

REM Build the Go application
go build -o "%tmpFile%" app\server.go

REM Execute the temporary file
"%tmpFile%"

REM Clean up the temporary file
del "%tmpFile%"