# Hardware

The kit is two boards: a small **core board** (the "yellow board", orange in practice, silkscreen "iCubeSmart") with the microcontroller, plugged onto the **base board** that carries the LEDs, drivers, keys and sensors.

## Microcontroller

| Item | Value |
|---|---|
| Chip | GigaDevice **GD32F103RET6**, LQFP64. Register-compatible with the STM32F103 (high density) |
| Core | ARM Cortex-M3, 72 MHz (8 MHz crystal × 9 PLL) |
| Flash | 512 KB at `0x08000000`, **2 KB pages** |
| SRAM | 64 KB at `0x20000000` |
| Device ID | `DBGMCU_IDCODE` (`0xE0042000`) reads `0x13090414`: DEV_ID `0x414` = high density |
| Flash size register | `0x1FFFF7E0` (16-bit, KB) |
| Unique ID | `0x1FFFF7E8` (96 bits) |
| Protection | Read protection level 0 (`0xA5`), no write protection |
| Bootloader | ROM UART bootloader on USART1 (PA9/PA10), 8E1. Entered with BOOT0 = 1 at reset |

Clocks: SYSCLK 72 MHz, APB2 72 MHz, APB1 36 MHz (timers on APB1 run at 72 MHz), ADC 12 MHz.

The core board footprint is drawn for an STM32F103RCT6; the fitted GD32F103RET6 has more flash and RAM. Datasheets: [GD32F103xx datasheet](../manuals/GD32F103xx-Datasheet-Rev-2.7.pdf), [GD32F10x user manual](../manuals/GD32F10x_User_Manual_EN_V2.3-3.pdf), [RM0008](../manuals/rm0008-stm32f101xx-stm32f102xx-stm32f103xx-stm32f105xx-and-stm32f107xx-advanced-armbased-32bit-mcus-stmicroelectronics.pdf).

## Core board

Schematic: [Yellow small board schematic](../diagrams/Yellow%20small%20board%20schematic%20黄色小板原理图.pdf) ("LETOO_RGB-CUBE_CORE V0.1").

| Part | Details |
|---|---|
| U4 AMS1117 | 5 V → 3.3 V regulator |
| Y2 | 8 MHz crystal |
| Reset button (J8) | NRST has a 10 kΩ pull-up and 1 µF to GND |
| BOOT0 (J2) | 2-pin header with a 10 kΩ resistor. A jumper selects the ROM bootloader at reset |
| BOOT1 | Unpopulated, so BOOT1 (PB2) = 0 |
| Serial header | 4-pin right angle: **GND, TXD, RXD, +5V**. USART1: TX = PA9, RX = PA10 |
| D1 / D2 | Status LEDs, green on **PB9**, red on **PB8**. Wired from 3.3 V through 1 kΩ to the pin, so **active low** |
| D3 | Power LED |
| 24C02 EEPROM | 256 bytes on I²C1: SCL = **PB6**, SDA = **PB7**, 10 kΩ pull-ups, device address `0x50` (`0xA0` write / `0xA1` read) |
| SD card | Drawn on the schematic (SDIO PC8–PC12, PD2, card detect PB5) but **not fitted** |
| Debug | No SWD header. SWDIO (PA13) and SWCLK (PA14) only reach the base connector, where they are KEY6 and KEY5 |

### Board-to-base connector (U26)

| Pin | Net | MCU | | Pin | Net | MCU |
|---|---|---|---|---|---|---|
| 1 | AIN1 | PA1 | | 20 | AIN0 | PA0 |
| 2 | PC3 | PC3 | | 19 | PC2 | PC2 |
| 3 | PC1 | PC1 | | 18 | PC0 | PC0 |
| 4 | RX2 | PA3 | | 17 | TX2 | PA2 |
| 5 | DB5V | 5 V in | | 16 | GND | |
| 6 | RST | NRST | | 15 | PB3 | PB3 |
| 7 | PA15 | PA15 | | 14 | PA14 | PA14 |
| 8 | USB_DP | PA12 | | 13 | DD1 | PA13 |
| 9 | 154EN | PA8 | | 12 | USB_DM | PA11 |
| 10 | PC6 | PC6 | | 11 | PC7 | PC7 |
| 21 | SPI1_CS | PA4 | | 26 | PC5 | PC5 |
| 22 | SPI1_CLK | PA5 | | 27 | GA | PB0 |
| 23 | SPI1_MISO | PA6 | | 28 | GB | PB1 |
| 24 | SPI1_MOSI | PA7 | | 29 | GC | PB10 |
| 25 | PC4 | PC4 | | 30 | GD | PB11 |

The `USB_DP`/`USB_DM` names are only net names; on this kit those pins are SW4 and KEY7.

## Base board

Silkscreen "iCubeSmart 8x8x RGB LED CUBE" (with tick boxes for 8x8x8 / 8x8x16 / 8x8x32). No schematic matches it exactly:

- The **pin functions** match [3D16RGB](../diagrams/3D16RGB%20Schematic%20diagram.pdf).
- The **LED driver count** (12) and the 32 layer pads match [8X8X8RGB](../diagrams/8X8X8RGB%20Schematic%20diagram.pdf). That sheet has no KEY7; this kit does.

| Part | Details |
|---|---|
| 12 × MBI5024 | 16-channel constant-current sinks for the LED cathodes: 192 outputs = 64 columns × RGB. See [display.md](display.md) |
| 74HC154 (×2 footprints) | 4-to-16 decoders selecting the layer (anode). The first (enable `154EN` = PA8) covers layers 1–16, the second (`154EN2` = PC7) layers 17–32. Only layers 1–8 are fitted |
| APM4953 | Dual P-MOSFETs switching 5 V onto each layer |
| 74HC245 (U27) | 3.3 V → 5 V buffer for SPI clock/data, LE and OE |
| KEY1–KEY7, RESET | Push buttons to GND |
| SWR1 | 4-way DIP switch, only SW1/SW2 are reachable |
| IR receiver | IRM-8601M2 on PC6 |
| Microphone | Electret mic M1 into an LM358 amplifier, into PA0 |
| J1 / J3 | 3.5 mm "cascade" jacks: RX (PA3) and TX (PA2) |
| J22 | 4-pin header: 5 V, TX2, RX2, GND (same UART2) |
| USB Type-C | **Power only** (5 V, no data). Pads are marked "8x8x8 5V 3A", "8x8x32 5V 7A" |
| TEST LED header | Unused by the firmware |

### Keys

All keys short to GND and have **no external pull-up**: the internal pull-up must be enabled. A press reads low.

| Key | MCU | Silkscreen | Default function |
|---|---|---|---|
| KEY1 | PC0 | Previous | Previous cartoon |
| KEY2 | PC1 | Next | Next cartoon |
| KEY3 | PC2 | Speed+ | Speed up |
| KEY4 | PB3 | Speed− | Speed down |
| KEY5 | PA14 | Run/Pause | Pause / resume |
| KEY6 | PA13 | Cycle | Play model (all cartoons / repeat one) |
| KEY7 | PA11 | ON/OFF | Display off / on |
| RESET | NRST | | Hardware reset |

The functions come from the key table of the running show ([show-engine.md](show-engine.md#key-functions)), so uploaded shows can change them.

PB3, PA13, PA14 and PA15 are JTAG/SWD pins after reset. The firmware must set `AFIO_MAPR.SWJ_CFG = 100` (debug port fully off) before KEY4, KEY5, KEY6 and SW3 work. See [firmware.md](firmware.md#mcu-pitfalls).

### DIP switches

Each switch has a 4.7 kΩ pull-up to 3.3 V; **ON connects to GND** (reads 0).

| Switch | MCU | Use |
|---|---|---|
| SW1 | PA1 | Mode bit |
| SW2 | PC3 | Mode bit |
| SW3 | PA15 | Not reachable on this base, reads 0. Only reported at startup (`Config2`) |
| SW4 | PA12 | Not reachable on this base, reads 0. Only reported at startup (`Config2`) |

Mode selection (1 = OFF/high, 0 = ON/low), as printed on the base:

| SW1 | SW2 | Mode | Number |
|---|---|---|---|
| 1 | 1 | Dynamic effects (default) | 0 |
| 0 | 1 | Audio effects | 1 |
| 1 | 0 | Test LED | 2 |
| 0 | 0 | Text effects (uploaded show) | 3 |

The factory firmware reads the switches only at reset. This firmware follows them while running ([firmware.md](firmware.md#modes)).

### IR receiver

IRM-8601M2 demodulating receiver on **PC6** with a 22 kΩ pull-up (10 kΩ on some revisions). Output is low during a 38 kHz burst. The kit's 21-key remote uses NEC coding; codes are listed in [serial-protocol.md](serial-protocol.md#ir-remote).

### Microphone

The LM358 output is biased around mid-scale and goes to **PA0** (ADC1 channel 0). With no sound the 12-bit reading sits near 1050; loud sound swings it towards 0 and 4095. See [show-engine.md](show-engine.md#audio-level).

### EEPROM contents

| Address | Size | Content |
|---|---|---|
| `0x00` | 8 | Player settings of this firmware ([storage.md](storage.md#eeprom)) |
| `0xA0` | 8 | Serial number (binary), written by `SetSerial` |
| `0xA8` | 8 | Cube type string, NUL padded, written by `SetType` |

Known cube type strings: `8x8x8A`, `8x8x8B`, `8x16x8A`, `161616A`, `161616B`, `161616C`, `8x32x8A`, `8x32x8B` (compared case-insensitively, index 0–7). This kit stores **`8x32x8b`**.

## Power

The cube runs from 5 V on the base's USB Type-C (up to 3 A for 8×8×8). The serial header's +5V pin can power the core board alone; do not connect it to an adapter's 5 V while Type-C is plugged in.

## Pin summary

| Function | Pins | Peripheral |
|---|---|---|
| Console, uploads, Pico, ROM bootloader | PA9 TX, PA10 RX | USART1 |
| Cascade jacks | PA2 TX, PA3 RX | USART2 |
| LED data | PA5 CLK, PA7 MOSI (PA6 MISO = chain output, unused) | SPI1 |
| LED latch / output enable | PC4 LE, PC5 OE (high = off) | GPIO |
| Layer address | PB0 A, PB1 B, PB10 C, PB11 D | GPIO |
| Layer enable | PA8 (layers 1–16), PC7 (17–32), low = on | GPIO |
| Keys | PC0, PC1, PC2, PB3, PA14, PA13, PA11 | GPIO + EXTI |
| DIP | PA1, PC3, PA15, PA12 | GPIO (polled) |
| IR | PC6 | GPIO + EXTI6 |
| Microphone | PA0 | ADC1 channel 0 |
| EEPROM | PB6 SCL, PB7 SDA | I2C1, 100 kHz |
| Status LEDs | PB9 green, PB8 red (active low) | GPIO |
