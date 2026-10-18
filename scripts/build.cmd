@echo off
rem Usage: scripts\build.cmd [package]   (default: ./cmd/8x8_rgb_cube)
setlocal
cd /d "%~dp0.."
set "PKG=%~1"
if "%PKG%"=="" set "PKG=./cmd/8x8_rgb_cube"
if not exist bin mkdir bin
echo Building %PKG% ...
tinygo build -o bin\main.hex -target targets\gd32f103re.json -opt s -panic print -gc conservative -size short "%PKG%"
if errorlevel 1 exit /b 1
echo Built bin\main.hex
