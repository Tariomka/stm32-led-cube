@echo off
rem Compares the Go show engine with the vendor Simulator.dll (needs the 32-bit Go toolchain target).
setlocal
cd /d "%~dp0.."
set GOARCH=386
go test -count=1 %* ./internal/lightshow/oracle/
