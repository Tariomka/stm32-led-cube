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
│   ├── documentation # Hardware, display, protocols, show format/engine, storage and firmware reference
│   ├── diagrams # Schematics and other wiring related docs
│   └── manuals # Data sheets, user manuals and similar docs
├── firmware # External binaries, for example dumped bootloader, etc.
├── internal # Core project functionality, considered as private library code
│   ├── runner # Entry point: creates the controllers and runs the main loop
│   ├── controller # One controller per concern, each owning its components:
│   │          # System (watchdog, LEDs, resets), Settings (EEPROM), Display, Keys, Remote (IR),
│   │          # Mode (DIP switches), Cascade (chained cubes), Show (player), Terminal (UART1 console + uploads)
│   ├── lightshow
│   │   └── builtin # Built-in light shows (the original cartoons as Go functions, played by the led-common-lib show engine)
│   ├── component # Hardware drivers: pins, buttons, layer decoder, display scanner (TIM2), microphone (ADC + TIM3), IR receiver, UARTs, EEPROM, flash, watchdog, LEDs
│   └── common # Small shared helpers
├── targets # Custom TinyGo target for the GD32F103RE (512 KB flash, 64 KB RAM)
├── tests # Host tests: the built-in shows render exactly like the factory show data
└── main.go # Windows task runner: go run . build/flash/dump/terminal/monitor/upload/test
```

## Light shows

The premade shows of the original firmware are Go functions in `internal/lightshow/builtin` (one per cartoon, generated once from the original show data; a test checks every frame matches the original). They run on a Go port of the original show engine, which lives in the [common library](https://github.com/Tariomka/led-common-lib) (`pkg/lightshow`) together with the UART protocol helpers.

DIP SW1/SW2 select the show mode (`controller.Mode`). Changing them switches the show right away, without a reset:

| SW1 | SW2 | Mode |
|---|---|---|
| 1 | 1 | Dynamic: the 41 original animations |
| 0 | 1 | Audio: 10 cartoons following the microphone |
| 1 | 0 | Test: LED test patterns |
| 0 | 0 | Text: shows uploaded from the PC (the original scrolling-text cartoon is not included) |

Uploading: make a show in the iCubeSmart "3D Software - RGB" and export `Update.hex`, then run `go run . upload [COMx]`, which starts the vendor "Download Software" at 115200 baud, and download `Update.hex`. The cube keeps only the show data (in flash, kept across power cycles) and plays all its cartoons in Text mode.

Other original features:
- KEY1..KEY7 run the original key functions (prev, next, speed+, speed-, pause, play mode, display off).
- Audio group cartoons follow the microphone level (PA0).
- The kit's 21-key IR remote works like the original: CH- display off, CH+ play mode, VOL-/EQ previous/next, VOL+ pause, >>| speed up, 100+ speed down, >|| speed reset, 200+ replay, 0..9 jump to cartoon 0..9.
- Speed, play mode, pause, display off and the selected cartoon are kept in the 24C02 EEPROM across power cycles.
- The original UART console works on UART1 (115200 baud): send `#?Help?!` (the vendor PC tool adds `#?`/`?!` itself; so does `go run . terminal`, an interactive terminal for Windows). Commands: Help, ReBoot, ReadMem, SetReg, GetReg, RegList, DelReg, ClrReg, SetType, SetSerial, Test. The registry is stored in the last 4 KB of flash.
- Remote control commands on the same console: `Key 1-7`, `Action 1-19`, `Cartoon <index>`, `Status`. The Raspberry Pi Pico uses this console and the upload frames, like the PC tools; there is no separate Pico protocol.
- Flash map: firmware up to 350 KB, uploaded show data at `0x08057800` (158 KB), registry at `0x0807F000` (see `targets/gd32f103re.ld`).
- Cascade: cubes chained over the 3.5 mm jacks (TX to the next cube's RX) find their order at power-up and play in step, following the first cube's keys. It uses the original protocol.

Full reference (hardware, display, protocols, show format, storage, firmware): [docs/documentation](docs/documentation/README.md).

Run the host tests with `go run . test` (or `go test ./tests/...`). `go run .` lists all tasks of the Windows task runner.

# Getting started

## Requirements

Bare minimum
1. Go v1.27.2 or above: [download link](https://go.dev/doc/install).
2. TinyGo v0.43 or above: [installation instructions](https://tinygo.org/getting-started/install/).
3. Any way to flash binaries to the microcontroller. I use STM32CubeProgrammer (listed in [Recommendations](#recommendations) section), but you can use other alternatives, example JTAG Programmer and `tinygo flash` command.

## Recommendations

1. Any text/code editor, example VS Code: [download link](https://code.visualstudio.com/download).
2. GNU Make.
3. STM32CubeProgrammer for flashing GD32 chip: [download link](https://www.st.com/en/development-tools/stm32cubeprog.html).

## Disclamer and Preparation

Unfinished. WORK IN PROGRESS! Do not use this.

## Quickstart

1. Build project binary (run from the repository root, the target file uses a relative linker script path):
    ```sh
    go run . build
    ```
    or with Make / directly:
    ```sh
    make build_version
    tinygo build -o bin/main.hex -target targets/gd32f103re.json -opt s -panic print -gc conservative ./cmd/8x8_rgb_cube
    ```
    > The default `bluepill-clone` target is limited to 64 KB of flash, which is too small for the original shows.

    On Windows `go run . flash [COMx]` builds and flashes over UART with the STM32CubeProgrammer CLI (BOOT0 jumper set), skipping step 2.
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
