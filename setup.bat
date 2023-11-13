@echo off
REM Author: @Benhexie
cls
echo :) Building...
cd src
go build -ldflags -H=windowsgui -o ../build/
cls
echo Build success :)
echo smile.exe is in build/
echo.
pause