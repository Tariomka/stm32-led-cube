@echo off
rem Usage: scripts\terminal.cmd [COMx] [baud] [command]   (defaults: COM4, 115200, interactive)
rem   scripts\terminal.cmd                    interactive console, lines are sent as #?line?!
rem   scripts\terminal.cmd COM4 115200 Help   send one command and print the reply
rem   scripts\terminal.cmd COM4 921600        the original firmware's console speed
setlocal
set "PORT=%~1"
if "%PORT%"=="" set "PORT=COM4"
set "BAUD=%~2"
if "%BAUD%"=="" set "BAUD=115200"
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0terminal.ps1" -Port %PORT% -Baud %BAUD% -Command "%~3"
