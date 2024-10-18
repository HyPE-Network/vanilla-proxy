@echo off

:: Install the new packages
go get -u all

:: Update packages used on your project
go mod tidy

:: Run the project
go run main.go