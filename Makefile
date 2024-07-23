project_build: create
	@echo "Starting to compile Tinygo binary, please wait..."
	@tinygo build -o ./build/main.elf -target=bluepill-clone ./src/main.go
	@echo "Build finished."

project_build_versioned: create
	@echo "Starting to compile versioned Tinygo binary, please wait..."
	@tinygo build -o ./build/main_$(shell date '+%Y_%m_%d_%H:%M').hex -target=bluepill-clone -size full ./src/main.go
	@echo "Build finished."

create:
	@mkdir -p build

clean:
	@rm -rf build
