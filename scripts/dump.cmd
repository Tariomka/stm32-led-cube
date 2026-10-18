@echo off
rem Usage: scripts\dump.cmd [COMx]   (default: COM4). Reads the whole 512 KB flash into firmware\dump_<timestamp>.bin
rem Needs the BOOT0 jumper fitted and a Reset before running.
setlocal
cd /d "%~dp0.."
set "PORT=%~1"
if "%PORT%"=="" set "PORT=COM4"
if "%CUBEPROG%"=="" set "CUBEPROG=C:\Program Files\STMicroelectronics\STM32Cube\STM32CubeProgrammer\bin\STM32_Programmer_CLI.exe"
for /f %%t in ('powershell -NoProfile -Command "Get-Date -Format yyyyMMdd_HHmmss"') do set "STAMP=%%t"
"%CUBEPROG%" -c port=%PORT% br=115200 -u 0x08000000 0x80000 "firmware\dump_%STAMP%.bin"
