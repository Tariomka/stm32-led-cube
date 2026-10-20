# Storage

## Flash map

Internal flash: 512 KB at `0x08000000`, 2 KB pages. An erased byte reads `0xFF`; programming writes 16-bit halfwords into erased cells.

| Range | Size | Content |
|---|---|---|
| `0x08000000–0x080577FF` | 350 KB | Firmware (linker script `targets/gd32f103re.ld`) |
| `0x08057800–0x0807EFFF` | 158 KB (79 pages) | **Show slot**: uploaded show data ([show-data.md](show-data.md)) |
| `0x0807F000–0x0807F7FF` | 2 KB | Registry |
| `0x0807F800–0x0807FFFF` | 2 KB | Registry backup |

- Show slot offset 0 holds the byte the upload addresses as `0x0800B000` (the show region of the vendor firmware image). The slot is larger than the 160 000-byte region, so any vendor show fits.
- An upload erases each slot page the first time it writes to it. Pages past the end of a smaller new show keep old bytes; they are never read because the show header gives the section sizes.
- The firmware reads the slot in place (no copy in RAM). An empty or invalid slot (no magic `0x12348765`) means "no uploaded show".
- Flash erase and program stall the CPU, so the display blanks briefly while the registry or the show slot is written.

### Factory firmware layout (for reference)

| Range | Content |
|---|---|
| `0x08000000–0x080045FF` | IAP bootloader (firmware update from SD card or UART, `IAP version: 2.10`) |
| `0x08005000`, `0x08005800` | Registry and backup |
| `0x08006000–0x0800AFFF` | Application code |
| `0x0800B000–0x080320FF` | Show data region (factory show: 32 100 bytes) |
| `0x08032100–0x08038FC3` | Drivers, startup code, more application code |

Flashing this firmware replaces all of it. Keep a full dump (`go run . dump`) to restore the factory firmware ([firmware.md](firmware.md#flashing)).

## Registry

A small key/value store for the console commands `SetReg`, `GetReg`, `RegList`, `DelReg`, `ClrReg`. The firmware itself never reads it; it is free for PC tools and users. The format is the factory firmware's, so an image is interchangeable.

Image (one 2 KB page, written to the registry page and then to the backup page; at startup the first valid one is loaded):

| Offset | Size | Content |
|---|---|---|
| 0 | 4 | u32 count of entries |
| 4 | 28 | `"Registry for 2.00"`, NUL padded |
| 32 | 16 × count | Entries: `key[12]` (NUL padded), u32 value size |
| 32 + 16·count | … | Values in entry order, each NUL terminated and padded to a multiple of 4 (the entry's size) |

Limits: key 1–11 characters, value 1–251 characters (value + NUL rounded up to 4 ≤ 256), at most 64 entries, whole image ≤ 2048 bytes. Keys are compared exactly. Values can't contain spaces (the console splits arguments on spaces).

If neither page holds a valid image the firmware prints `Regedit format error!` at startup and starts with an empty registry.

## EEPROM

24C02: 256 bytes, 8-byte write pages, I²C address `0x50` on I2C1 (PB6 SCL, PB7 SDA, 100 kHz). A write waits for the chip's internal write cycle (≤ 5 ms) by polling its address. A failed transfer is repeated once after resetting the I2C peripheral (the STM32F1 I2C peripheral can get stuck busy, and TinyGo's driver doesn't guard its read sequence against interrupts).

| Address | Size | Content | Written by |
|---|---|---|---|
| `0x00–0x07` | 8 | Player settings (below) | This firmware, 2 s after the last change |
| `0xA0–0xA7` | 8 | Serial number | `SetSerial` |
| `0xA8–0xAF` | 8 | Cube type string, NUL padded (this kit: `8x32x8b`) | `SetType` |

The factory firmware checks the serial number against a digest and stops if it doesn't match; this firmware only prints the type (`T='…'`) and never checks the serial.

### Settings record

| Byte | Content |
|---|---|
| 0, 1 | Magic `'L'`, `'C'` |
| 2 | Version, 1 |
| 3 | Speed index 0–7 |
| 4 | Flags: bit 0 play all, bit 1 paused, bit 2 display off |
| 5 | Mode the cartoon belongs to (0 Dynamic, 1 Audio, 2 Test, 3 Text) |
| 6 | Cartoon index |
| 7 | Checksum: bitwise NOT of the 8-bit sum of bytes 0–6 |

- Saved 2 s after the last key, IR, console action or mode change, so a burst of presses writes once.
- At power-up a valid record restores speed, play model, display off and pause (a paused cube first draws one frame, so it isn't blank). The cartoon is restored only if the DIP switches select the same mode as saved.
- A blank or foreign EEPROM (bad magic, version or checksum) leaves the defaults.
