@echo off
go tool pack grc %1.8.8 %1.8 %1.rc.o
go tool 8l -L %GOPATH%\pkg\windows_386 -o %1.exe -s  %1.8.8