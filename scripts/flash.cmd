@echo off
rem Usage: scripts\flash.cmd [COMx] [file]   (defaults: COM4, bin\main.hex)
rem Needs the BOOT0 jumper fitted and a Reset before running. Builds first when flashing bin\main.hex.
setlocal
cd /d "%~dp0.."
set "PORT=%~1"
if "%PORT%"=="" set "PORT=COM4"
set "FILE=%~2"
if "%FILE%"=="" (
    set "FILE=bin\main.hex"
    call scripts\build.cmd
    if errorlevel 1 exit /b 1
)
if "%CUBEPROG%"=="" set "CUBEPROG=C:\Program Files\STMicroelectronics\STM32Cube\STM32CubeProgrammer\bin\STM32_Programmer_CLI.exe"
"%CUBEPROG%" -c port=%PORT% br=115200 -w "%FILE%" -v -g 0x08000000
