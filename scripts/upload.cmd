@echo off
rem Usage: scripts\upload.cmd [COMx]   (default COM4)
rem Starts the vendor "Download Software" at 115200 baud (our UART1 speed). In the tool, pick the
rem Update.hex made by "3D Software - RGB" and start the download: the cube keeps the show data
rem and plays it in Text mode (DIP SW1 = 0, SW2 = 0). Set ICUBE_DOWNLOAD to use another copy.
setlocal
cd /d "%~dp0.."
set "PORT=%~1"
if "%PORT%"=="" set "PORT=COM4"
if "%ICUBE_DOWNLOAD%"=="" set "ICUBE_DOWNLOAD=.agent\external\3D Software - RGB Cube\Download Software.exe"
start "" "%ICUBE_DOWNLOAD%" %PORT% 115200
