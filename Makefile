build: create
	@echo "Starting to compile Tinygo binary, please wait..."
	@tinygo build -o ./bin/main.elf -target=bluepill-clone ./src/main.go
	@echo "Build finished."

build_version: create
	@echo "Starting to compile versioned Tinygo binary, please wait..."
	@tinygo build -o ./bin/main_$(shell date '+%Y_%m_%d_%H:%M').hex -target=bluepill-clone -size full ./src/main.go
	@echo "Build finished."

create:
	@mkdir -p bin

clean:
	@rm -rf bin

test:
	@tinygo test -target=bluepill-clone ./test/common/error_test.go
