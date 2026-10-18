@echo off
rem Usage: scripts\monitor.cmd [COMx] [baud]   (defaults: COM4, 115200)
setlocal
set "PORT=%~1"
if "%PORT%"=="" set "PORT=COM4"
set "BAUD=%~2"
if "%BAUD%"=="" set "BAUD=115200"
tinygo monitor -port %PORT% -baudrate %BAUD%
