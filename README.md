# STM32 Led Cube - 3D RGB LED cube project

This is a fun little personal open-source project for a [iCubeSmart](https://icubesmart.com/) 3D8RGB Led Cube Kit, which may or may not be found [here](https://icubesmart.com/products/icubesmart-3d8rgb-led-cube-kit-full-color-8x8x8-cube-diy-electronic-kits).
It was inspired by [Sliicy's](https://github.com/Sliicy/8x8x8-LED) and [andypants152's](https://github.com/andypants152/iCubeSmart) projects and written in Golang (TinyGo compiler) for me to learn something new.

The brains of the operation is GD32F103RET6 microcontroller (`Bluepill clone`) and as such the project is designed, compiled and ran on this microcontroller.

This is one part of a bigger 3 part project I'm cooking:
1. The main worker - STM32/GD32 board (this project)
2. A communication middle man for wireless control - Raspberry Pico W using [this project](https://github.com/Tariomka/rpi-led-communication)
3. Desktop app for designing light shows and monitoring both microcontrollers' activity - [this project](https://github.com/Tariomka/desktop-led-controller)

I'm also using a [common library](https://github.com/Tariomka/led-common-lib) for all 3 projects to not have duplicate and/or out of sync logic components.

# Project Structure

Here's a basic brakedown of the structure of directories:
```bash
├── .vscode # Editor related files
├── bin # Compiled binary files
├── cmd # Applications directory
│   └── **/main.go # Main entry point for eacy application
├── docs # Documentation directory
│   ├── diagrams # Schematics and other wiring related docs
│   └── manuals # Data sheets, user manuals and similar docs
├── firmware # External binaries, for example dumped bootloader, etc.
└── internal # Core project functionality, considered as private library code
```


# Requirements

1. Go v1.23.4 or above: [download link](https://go.dev/doc/install).
2. TinyGo v0.35 or above: [installation instructions](https://tinygo.org/getting-started/install/).
3. Any code editor, ex. VS Code.
4. GNU Make.
5. STM32CubeProgrammer for flashing GD32 chip: [download link](https://www.st.com/en/development-tools/stm32cubeprog.html).

# Getting started

