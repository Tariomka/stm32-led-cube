# STM32 Led Cube - 3D RGB LED cube project

This is a fun little personal open-source project for a [iCubeSmart](https://icubesmart.com/) 3D8RGB Led Cube Kit, which may or may not be found [here](https://icubesmart.com/products/icubesmart-3d8rgb-led-cube-kit-full-color-8x8x8-cube-diy-electronic-kits).
It was inspired by [Sliicy's](https://github.com/Sliicy/8x8x8-LED) and [andypants152's](https://github.com/andypants152/iCubeSmart) projects and written in Golang (TinyGo compiler) for me to learn something new.

The brains of the operation is GD32F103RET6 microcontroller (`Bluepill clone`) and as such the project is designed, compiled and ran on this microcontroller.

This is one part of a bigger 3 part project I'm cooking:
1. The main worker - STM32/GD32 board (this project)
2. A communication middle man for wireless control - Raspberry Pico W using [this project](https://github.com/Tariomka/rpi-led-communication)
3. Desktop app for designing light shows and monitoring both microcontrollers' activity - [this project](https://github.com/Tariomka/desktop-led-controller)

I'm also using a [common library](https://github.com/Tariomka/led-common-lib) for all 3 projects to not have duplicate and/or out of sync logic components.

<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#stm32-led-cube---3d-rgb-led-cube-project">About The Project</a></li>
    <li><a href="#project-structure">Project Structure</a></li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#requirements">Requirements</a></li>
        <li><a href="#recommendations">Recommendations</a></li>
      </ul>
    </li>
  </ol>
</details>

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

# Getting started

## Requirements

Bare minimum
1. Go v1.24.6 or above: [download link](https://go.dev/doc/install).
2. TinyGo v0.39 or above: [installation instructions](https://tinygo.org/getting-started/install/).
3. Any way to flash binaries to the microcontroller. I use STM32CubeProgrammer (listed in [Recommendations](#recommendations) section), but you can use other alternatives, example JTAG Programmer and `tinygo flash` command.

## Recommendations

1. Any text/code editor, example VS Code: [download link](https://code.visualstudio.com/download).
2. GNU Make.
3. STM32CubeProgrammer for flashing GD32 chip: [download link](https://www.st.com/en/development-tools/stm32cubeprog.html).

## Disclamer and Preparation

Unfinished. WORK IN PROGRESS! Do not use this.

## Quickstart

1. Build project binary with Make:
    ```sh
    make
    ```
    or without Make:
    ```sh
    tinygo build -target bluepill-clone ./cmd/8x8_rgb_cube/main.go
    ```
2. Flash newly created binary:
    <!-- TODO: Add screenshots -->
    1. Short the `BOOT0` pins on the GD32 yellow board
    2. Plug in GD32 to your computer via UART programmer:
        | GD32 pins | UART Programmer pins |
        | :-------: | :------------------: |
        | GND       | GND                  |
        | TXD       | TXD                  |
        | RXD       | RXD                  |
        | +5V       | 5V                   |
        |           |    Short VCC - 3V3   |
    3. Launch STM32CubeProgrammer
    4. In the dropdown select `UART`
    5. Select Port that corresponds to the connected GD32, typically its `ttyUSB0` on Linux
    6. Click `Connect`
    > NOTE: if you get an error like this:
    >
    > > Error: Activating device: KO. Please, verify the boot mode configuration and check the serial port configuration. Reset your device then try again...
    >
    > double check steps 2.I and 2.V and if everything looks correct, press the `Reset` button on the Yellow Board and step 6 again.
    7. Navigate to `Erasing & programming` section
    8. Click `Browse` and find the binary that was compiled in step 1
    > If the project was compiled using Make, the binary will be located in `./bin/main_8x8_rgb_{date_time}.hex`, otherwise it should be localed in the root directory named `main.elf`
    9. Click `Start Programming` and wait until upload finishes.
