#!/bin/bash
echo "Starting to compile Tinygo binary, please wait..."

mkdir -p build
current_date_time="`date +%Y_%m_%d_%H:%M`";
tinygo build -o ./build/main_${current_date_time}.hex -target bluepill-clone -size full ./src/main.go

echo "Build finished."