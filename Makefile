BIN_DIR = bin
ifdef OS
	VERSION = $(strip $(shell cmd /C date /t))_$(subst :,-,$(shell cmd /C time /t))
	RM = del /s /q
else
	VERSION = $(shell date '+%Y_%m_%d_%H:%M')
	RM = rm -rf
endif

build: create
	@echo "Starting to compile Tinygo binary, please wait..."
	@tinygo build -o ./$(BIN_DIR)/main.elf -target=bluepill-clone ./src/main.go
	@echo "Build finished."

build_version: create
	@echo "Starting to compile versioned Tinygo binary, please wait..."
	@tinygo build -o ./$(BIN_DIR)/main_$(VERSION).hex -target=bluepill-clone -size full ./src/main.go
	@echo "Build finished."

create:
	@if not exist $(BIN_DIR) mkdir $(BIN_DIR)

clean:
	@$(RM) $(BIN_DIR)

test:
	@tinygo test -target=bluepill-clone ./test/common/error_test.go
