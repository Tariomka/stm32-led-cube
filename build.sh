#!/bin/bash
echo "Starting to compile Tinygo binary, please wait..."

mkdir -p bin
current_date_time="`date +%Y_%m_%d_%H:%M`";
tinygo build -o ./bin/main_${current_date_time}.hex -target=bluepill-clone ./src/main.go

echo "Build finished."