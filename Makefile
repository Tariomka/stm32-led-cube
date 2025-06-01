BIN_DIR = bin
TINYGO_FLAGS = -opt s -panic print -gc conservative -size full
TARGET = -target bluepill-clone
ifdef OS
	VERSION = v$(strip $(shell cmd /C date /t))_$(subst :,-,$(shell cmd /C time /t))
	RM = del /s /q
else
	VERSION = $(shell date '+%Y_%m_%d_%H:%M')
	RM = rm -rf
endif
BIN_NAME = main_8x8_rgb_$(VERSION).hex

build_version: create
	@echo Starting to compile versioned Tinygo binary, please wait...
	@tinygo build -o ./$(BIN_DIR)/$(BIN_NAME) $(TARGET) $(TINYGO_FLAGS) ./cmd/8x8_rgb_cube/main.go
	@echo Created '$(BIN_NAME)' binary file.
	@echo Build finished.

build: create
	@echo Starting to compile Tinygo binary, please wait...
	@tinygo build -o ./$(BIN_DIR)/main_8x8_rgb.elf $(TARGET) -opt s ./cmd/8x8_rgb_cube/main.go
	@echo Build finished.

create:
	@if [ ! -d $(BIN_DIR) ]; then mkdir $(BIN_DIR); fi

clean:
	@$(RM) $(BIN_DIR)
